package calendar

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

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

type GetCalendarListData struct {
	CalendarList []models.Calendar
	AccessToken  string
}

// Calls GetCalendarList but broadcasts the result to channel
func GetCalendarListAsync(accessToken string, c chan GetCalendarListData) {
	calendars, _ := GetCalendarList(accessToken)
	c <- GetCalendarListData{calendars, accessToken}
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
		c <- nil
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

func GetUsersCalendarEvents(user *models.User, accounts models.Set[string], timeMin time.Time, timeMax time.Time) ([]models.CalendarEvent, *errs.GoogleAPIError) {
	db.RefreshUserTokenIfNecessary(user, accounts)

	// Map mapping access token to the calendar list associated with that access token
	calendarListMap := make(map[string][]models.Calendar)

	// Get primary user's calendar, throw error if gcal access not granted
	calendars, err := GetCalendarList(user.AccessToken)
	if err != nil {
		return nil, err
	}
	calendarListMap[user.AccessToken] = calendars

	// Get secondary account calendars
	calendarListChan := make(chan GetCalendarListData)
	for _, account := range user.CalendarAccounts {
		if !*account.Enabled {
			continue
		}

		go GetCalendarListAsync(account.AccessToken, calendarListChan)
	}
	for _, account := range user.CalendarAccounts {
		if !*account.Enabled {
			continue
		}

		calendarListData := <-calendarListChan
		if calendarListData.CalendarList != nil {
			calendarListMap[calendarListData.AccessToken] = calendarListData.CalendarList
		}
	}

	// Get a list of calendar events from the each gcal account associated with a user
	calendarEventsChan := make(chan []models.CalendarEvent)
	calendarEvents := make([]models.CalendarEvent, 0)
	for accessToken, calendarList := range calendarListMap {
		for _, calendar := range calendarList {
			go GetCalendarEventsAsync(accessToken, calendar.Id, timeMin, timeMax, calendarEventsChan)
		}
	}
	for _, calendarList := range calendarListMap {
		for range calendarList {
			events := <-calendarEventsChan
			if events != nil {
				calendarEvents = append(calendarEvents, events...)
			}
		}
	}

	return calendarEvents, nil
}
