package calendar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"schej.it/server/db"
	"schej.it/server/errs"
	"schej.it/server/logger"
	"schej.it/server/models"
)

// Get the user's list of calendars
func GetCalendarList(accessToken string) ([]models.Calendar, *errs.GoogleAPIError) {
	// TODO: update user object with calendars and allow for customization of whether or not to show calendar in schedule
	req, _ := http.NewRequest(
		"GET",
		"https://www.googleapis.com/calendar/v3/users/me/calendarList?fields=items(id,summary,selected)",
		nil,
	)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.StdErr.Panicln(err)
	}
	defer resp.Body.Close()

	// Define stucts to parse json response
	type Response struct {
		Items []models.Calendar   `json:"items"`
		Error errs.GoogleAPIError `json:"error"`
	}

	// Parse the response
	var res Response
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		logger.StdErr.Panicln(err)
	}

	// Check if the response returned an error
	if res.Error.Errors != nil {
		return nil, &res.Error
	}

	// Append only the selected calendars
	calendars := make([]models.Calendar, 0)
	for _, calendar := range res.Items {
		if calendar.Selected {
			calendars = append(calendars, calendar)
		}
	}

	return calendars, nil
}

// Calls GetCalendarList but broadcasts the result to channel
func GetCalendarListAsync(accessToken string, c chan []models.Calendar) {
	calendars, err := GetCalendarList(accessToken)
	if err == nil {
		c <- calendars
	}
}

// Get the user's list of calendar events for the given calendar
func GetCalendarEventsAsync(accessToken string, calendarId string, timeMin time.Time, timeMax time.Time, c chan []models.CalendarEvent) ([]models.CalendarEvent, *errs.GoogleAPIError) {
	min, _ := timeMin.MarshalText()
	max, _ := timeMax.MarshalText()
	//fmt.Printf("https://www.googleapis.com/calendar/v3/calendars/%s/events?timeMin=%s&timeMax=%s&singleEvents=true\n", url.PathEscape(calendarId), min, max)
	req, _ := http.NewRequest(
		"GET",
		fmt.Sprintf("https://www.googleapis.com/calendar/v3/calendars/%s/events?fields=items(summary,start,end)&timeMin=%s&timeMax=%s&singleEvents=true", url.PathEscape(calendarId), min, max),
		nil,
	)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.StdErr.Panicln(err)
	}
	defer resp.Body.Close()

	// Define some structs to parse the json response
	type TimeInfo struct {
		DateTime time.Time `json:"dateTime" binding:"required"`
	}
	type Item struct {
		Summary string   `json:"summary"`
		Start   TimeInfo `json:"start"`
		End     TimeInfo `json:"end"`
	}
	type Response struct {
		Items []Item              `json:"items"`
		Error errs.GoogleAPIError `json:"error"`
	}

	// Parse the response
	var res Response
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		logger.StdErr.Panicln(err)
	}

	// Check if the response returned an error
	if res.Error.Errors != nil {
		return nil, &res.Error
	}

	// Format response to return
	calendarEvents := make([]models.CalendarEvent, 0)
	for _, item := range res.Items {
		// Don't include events that are all day events
		// Don't include events that are greater than 24 hours
		if item.Start.DateTime.IsZero() || item.End.DateTime.Sub(item.Start.DateTime).Hours() >= 24 {
			continue
		}

		// Restructure event
		calendarEvents = append(calendarEvents, models.CalendarEvent{
			Summary:   item.Summary,
			StartDate: primitive.NewDateTimeFromTime(item.Start.DateTime),
			EndDate:   primitive.NewDateTimeFromTime(item.End.DateTime),
		})
	}

	c <- calendarEvents

	return calendarEvents, nil
}

func GetUsersCalendarEvents(user *models.User, timeMin time.Time, timeMax time.Time) ([]models.CalendarEvent, *errs.GoogleAPIError) {
	db.RefreshUserTokenIfNecessary(user)

	// Get primary user's calendar, throw error if gcal access not granted
	calendars, err := GetCalendarList(user.AccessToken)
	if err != nil {
		return nil, err
	}

	// Get secondary account calendars

	// Call the google calendar API to get a list of calendar events from the user's gcal
	calendarEventsChan := make(chan []models.CalendarEvent)
	calendarEvents := make([]models.CalendarEvent, 0)
	for _, calendar := range calendars {
		go GetCalendarEventsAsync(user.AccessToken, calendar.Id, timeMin, timeMax, calendarEventsChan)
	}
	for range calendars {
		events := <-calendarEventsChan
		calendarEvents = append(calendarEvents, events...)
	}

	return calendarEvents, nil
}

func ScheduleEvent(user *models.User, eventName string, eventId string, calendarEventId string, startDate primitive.DateTime, endDate primitive.DateTime, attendeeEmails []string, location string, description string) (*string, *errs.GoogleAPIError) {
	db.RefreshUserTokenIfNecessary(user)

	attendees := make(bson.A, 0)
	attendees = append(attendees, bson.M{"email": user.Email, "responseStatus": "accepted"})
	for _, email := range attendeeEmails {
		attendees = append(attendees, bson.M{"email": email, "responseStatus": "needsAction"})
	}

	body, _ := json.Marshal(bson.M{
		"start": bson.M{
			"dateTime": startDate,
		},
		"end": bson.M{
			"dateTime": endDate,
		},
		"attendees":   attendees,
		"summary":     eventName,
		"description": fmt.Sprintf("%s\n\nThis event was scheduled with schej: https://schej.it/e/%s", description, eventId),
		"location":    location,
	})
	reqBody := bytes.NewBuffer(body)

	// Create calendar event
	var req *http.Request
	if len(calendarEventId) > 0 {
		// Update existing event
		req, _ = http.NewRequest(
			"PUT",
			fmt.Sprintf("https://www.googleapis.com/calendar/v3/calendars/primary/events/%s?fields=id&sendUpdates=all", calendarEventId),
			reqBody,
		)
	} else {
		// Create new event
		req, _ = http.NewRequest(
			"POST",
			"https://www.googleapis.com/calendar/v3/calendars/primary/events?fields=id&sendUpdates=all",
			reqBody,
		)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", user.AccessToken))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.StdErr.Panicln(err)
	}
	defer resp.Body.Close()

	// Parse the response
	response := struct {
		Id    string              `json:"id"`
		Error errs.GoogleAPIError `json:"error"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		logger.StdErr.Panicln(err)
	}

	// Check if the response returned an error
	if response.Error.Errors != nil {
		return nil, &response.Error
	}

	return &response.Id, nil
}

func UnscheduleEvent(user *models.User, calendarEventId string) *errs.GoogleAPIError {
	db.RefreshUserTokenIfNecessary(user)

	req, _ := http.NewRequest(
		"DELETE",
		fmt.Sprintf("https://www.googleapis.com/calendar/v3/calendars/primary/events/%s?fields=id&sendUpdates=all", calendarEventId),
		nil,
	)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", user.AccessToken))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.StdErr.Panicln(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		// Parse the response
		response := struct {
			Error errs.GoogleAPIError `json:"error"`
		}{}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			logger.StdErr.Panicln(err)
		}

		// Check if the response returned an error
		if response.Error.Errors != nil {
			return &response.Error
		}
	}

	return nil
}
