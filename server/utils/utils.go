package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"

	"github.com/brianvoe/sjwt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"schej.it/server/errs"
	"schej.it/server/logger"
	"schej.it/server/models"
)

// Returns whether running on production server
func IsRelease() bool {
	mode := os.Getenv("GIN_MODE")
	return mode == "release"
}

func PrintJson(s gin.H) {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		logger.StdErr.Panicln(err)
	}

	fmt.Println(string(data))
}

func ParseJWT(jwt string) sjwt.Claims {
	claims, err := sjwt.Parse(jwt)
	if err != nil {
		logger.StdErr.Panicln(err)
	}

	return claims
}

func StringToObjectID(s string) primitive.ObjectID {
	objectID, err := primitive.ObjectIDFromHex(s)
	if err != nil {
		logger.StdErr.Panicln(err)
	}

	return objectID
}

// Gets the user id from the current session as an ObjectID object
func GetUserId(session sessions.Session) primitive.ObjectID {
	return StringToObjectID(session.Get("userId").(string))
}

// Gets the access token expire date from an "expiresIn" int representing the number of seconds
// after which the access token will expire
func GetAccessTokenExpireDate(expiresIn int) time.Time {
	expireDuration, err := time.ParseDuration(fmt.Sprintf("%ds", expiresIn))
	if err != nil {
		logger.StdErr.Panicln(err)
	}
	return time.Now().Add(expireDuration)
}

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

// Get the user's list of calendar events for the given calendar
func GetCalendarEvents(accessToken string, calendarId string, timeMin time.Time, timeMax time.Time) ([]models.CalendarEvent, *errs.GoogleAPIError) {
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
		// Restructure event
		calendarEvents = append(calendarEvents, models.CalendarEvent{
			Summary:   item.Summary,
			StartDate: primitive.NewDateTimeFromTime(item.Start.DateTime),
			EndDate:   primitive.NewDateTimeFromTime(item.End.DateTime),
		})
	}

	return calendarEvents, nil
}

// Returns the ISO date string for the given date
func GetDateString(date time.Time) string {
	s, _ := date.UTC().MarshalText()
	return string(s)[:10]
}

// Returns a time object with the given date and a time string in the form of "00:00:00"
func GetDateAtTime(date time.Time, timeString string) time.Time {
	utcDateString := GetDateString(date)
	newDate, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT%sZ", utcDateString, timeString))
	if err != nil {
		logger.StdErr.Panicln(err)
	}
	return newDate
}

// Inserts the given value at the specified index in the slice. Returns the updated slice
func Insert[T any](arr []T, index int, value T) ([]T, error) {
	if index < 0 {
		return nil, errors.New("index cannot be less than 0")
	}

	if index >= len(arr) {
		return append(arr, value), nil
	}

	arr = append(arr[:index+1], arr[index:]...)
	arr[index] = value

	return arr, nil
}

// Returns whether the given slice contains the given value
func Contains[T comparable](arr []T, value T) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}

// Escapes regex for a string
func EscapeRegExp(str string) string {
	check := regexp.MustCompile(`([.*+?^${}()|[\]\\])`)
	return check.ReplaceAllString(str, "\\${1}")
}
