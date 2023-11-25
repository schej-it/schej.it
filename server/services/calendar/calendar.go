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
	"schej.it/server/utils"
)

// Get the user's list of calendars
func GetCalendarList(accessToken string) (map[string]models.SubCalendar, *errs.GoogleAPIError) {
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
	calendars := make(map[string]models.SubCalendar)
	for _, calendar := range res.Items {
		var enabled *bool
		if calendar.Selected {
			enabled = utils.TruePtr()
		} else {
			enabled = utils.FalsePtr()
		}

		calendars[calendar.Id] = models.SubCalendar{
			Name:    calendar.Summary,
			Enabled: enabled,
		}
	}

	return calendars, nil
}

type GetCalendarListData struct {
	CalendarList map[string]models.SubCalendar `json:"calendarList"`
	AccessToken  string                        `json:"accessToken"`
	Email        string                        `json:"email"`
	Error        *errs.GoogleAPIError          `json:"error"`
}

// Calls GetCalendarList but broadcasts the result to channel
func GetCalendarListAsync(email string, accessToken string, c chan GetCalendarListData) {
	// Recover from panics
	defer func() {
		if err := recover(); err != nil {
			c <- GetCalendarListData{Error: &errs.GoogleAPIError{Message: err.(string)}}
		}
	}()

	calendars, err := GetCalendarList(accessToken)

	c <- GetCalendarListData{CalendarList: calendars, AccessToken: accessToken, Email: email, Error: err}
}

type GetCalendarEventsData struct {
	CalendarEvents []models.CalendarEvent `json:"calendarEvents"`
	Email          string                 `json:"email"`
	Error          *errs.GoogleAPIError   `json:"error"`
}

// Get the user's list of calendar events for the given calendar
func GetCalendarEventsAsync(email string, accessToken string, calendarId string, timeMin time.Time, timeMax time.Time, c chan GetCalendarEventsData) {
	min, _ := timeMin.MarshalText()
	max, _ := timeMax.MarshalText()
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
	type Response struct {
		Items []struct {
			Summary string `json:"summary"`
			Start   struct {
				DateTime time.Time `json:"dateTime" binding:"required"`
			} `json:"start"`
			End struct {
				DateTime time.Time `json:"dateTime" binding:"required"`
			} `json:"end"`
		} `json:"items"`
		Error *errs.GoogleAPIError `json:"error"`
	}

	// Parse the response
	var res Response
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		logger.StdErr.Panicln(err)
	}

	// Check if the response returned an error
	if res.Error != nil {
		c <- GetCalendarEventsData{Email: email, Error: res.Error}
		return
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

	c <- GetCalendarEventsData{CalendarEvents: calendarEvents, Email: email}
}

type CalendarEventsWithError struct {
	CalendarEvents []models.CalendarEvent `json:"calendarEvents"`
	Error          *errs.GoogleAPIError   `json:"error,omitempty"`
}

// Returns a map mapping email to the calendar events associated with that email, and an error if there was an error fetching events for that email
func GetUsersCalendarEvents(user *models.User, accounts models.Set[string], timeMin time.Time, timeMax time.Time) (map[string]CalendarEventsWithError, bool) {
	db.RefreshUserTokenIfNecessary(user, accounts)

	returnAllAccounts := len(accounts) == 0
	editedCalendarAccounts := false

	calendarEventsMap := make(map[string]CalendarEventsWithError)

	calendarListChan := make(chan GetCalendarListData)
	calendarEventsChan := make(chan GetCalendarEventsData)

	// Get calendar lists
	numCalendarListRequests := 0
	for _, account := range user.CalendarAccounts {
		// Get secondary account  calendars
		if _, ok := accounts[account.Email]; ok || returnAllAccounts {
			go GetCalendarListAsync(account.Email, account.AccessToken, calendarListChan)
			numCalendarListRequests++

			calendarEventsMap[account.Email] = CalendarEventsWithError{
				CalendarEvents: make([]models.CalendarEvent, 0),
			}
		}
	}

	// After each calendar list is fetched, get the calendar events from each calendar
	numCalendarEventsRequests := 0
	for i := 0; i < numCalendarListRequests; i++ {
		calendarListData := <-calendarListChan

		if calendarListData.Error != nil {
			// This is needed to be able to send an error back to user if a given calendar account's refresh token is invalid, for example
			go func() { // needs to be async because writing to a channel is blocking
				calendarEventsChan <- GetCalendarEventsData{Email: calendarListData.Email, Error: calendarListData.Error}
			}()
			numCalendarEventsRequests++
			continue
		}

		// Edit subcalendars map
		account := user.CalendarAccounts[calendarListData.Email]
		if account.SubCalendars == nil {
			account.SubCalendars = &calendarListData.CalendarList
			user.CalendarAccounts[calendarListData.Email] = account
			editedCalendarAccounts = true
		} else {
			// Add subCalendar if it doesn't exist
			for id, subCalendar := range calendarListData.CalendarList {
				if _, ok := (*account.SubCalendars)[id]; !ok {
					(*account.SubCalendars)[id] = subCalendar

					if !editedCalendarAccounts {
						editedCalendarAccounts = true
					}
				}
			}

			// Remove subCalendar if it no longer exists
			for id := range *account.SubCalendars {
				if _, ok := calendarListData.CalendarList[id]; !ok {
					delete(*account.SubCalendars, id)

					if !editedCalendarAccounts {
						editedCalendarAccounts = true
					}
				}
			}
		}
		user.CalendarAccounts[calendarListData.Email] = account

		for id, calendar := range *account.SubCalendars {
			if *calendar.Enabled {
				go GetCalendarEventsAsync(calendarListData.Email, calendarListData.AccessToken, id, timeMin, timeMax, calendarEventsChan)
				numCalendarEventsRequests++
			}
		}
	}

	// After calendar events are fetched, append to the calendarEvents array associated with the given email
	for i := 0; i < numCalendarEventsRequests; i++ {
		calendarEventsData := <-calendarEventsChan
		email := calendarEventsData.Email

		if _, ok := calendarEventsMap[email]; !ok {
			calendarEventsMap[email] = CalendarEventsWithError{}
		}

		if events, ok := calendarEventsMap[email]; ok {
			if calendarEventsData.Error != nil {
				events.Error = calendarEventsData.Error
			} else {
				events.CalendarEvents = append(events.CalendarEvents, calendarEventsData.CalendarEvents...)
			}
			calendarEventsMap[email] = events
		}
	}

	return calendarEventsMap, editedCalendarAccounts
}
