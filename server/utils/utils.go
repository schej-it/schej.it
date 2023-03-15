package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"

	"github.com/brianvoe/sjwt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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
	// defer resp.Body.close() ??
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

// Escapes regex for a string
func EscapeRegExp(str string) string {
	check := regexp.MustCompile(`([.*+?^${}()|[\]\\])`)
	return check.ReplaceAllString(str, "\\${1}")
}

// Returns the correct client id given the token origin
func GetClientIdFromTokenOrigin(tokenOrigin models.TokenOriginType) string {
	switch tokenOrigin {
	case models.ANDROID:
		return os.Getenv("ANDROID_CLIENT_ID")
	case models.IOS:
		return os.Getenv("IOS_CLIENT_ID")
	default:
		return os.Getenv("CLIENT_ID")
	}
}

// Prints the http response as a string
func PrintHttpResponse(resp *http.Response) {
	body, _ := ioutil.ReadAll(resp.Body)
	logger.StdOut.Println(string(body))
	resp.Body = io.NopCloser(bytes.NewBuffer(body))
}

func AddUserToMailchimp(email string, firstName string, lastName string) {
	// Adds the given user to the default mailchimp audience
	apiKey := os.Getenv("MAILCHIMP_API_KEY")

	body, _ := json.Marshal(bson.M{
		"email_address": email, "status": "subscribed", "merge_fields": bson.M{
			"FNAME": firstName,
			"LNAME": lastName,
		},
		"tags": bson.A{"user"},
	})
	bodyBuffer := bytes.NewBuffer(body)

	req, _ := http.NewRequest("POST", "https://us21.api.mailchimp.com/3.0/lists/b5c79106b4/members", bodyBuffer)
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", apiKey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.StdErr.Println(err)
	}
	defer resp.Body.Close()
}

func AddUserToMailjet(email string, firstName string, lastName string, picture string) {
	// Adds the given user to the Mailjet contact list
	apiKey := os.Getenv("MAILJET_API_KEY")
	apiSecret := os.Getenv("MAILJET_API_SECRET")
	listId := os.Getenv("MAILJET_LIST_ID")

	// Create contact
	// POST https://api.mailjet.com/v3/REST/contact {"EMAIL", email}
	// contactId = result.Data[0].ID (integer)
	body, _ := json.Marshal(bson.M{
		"Email": email,
	})
	bodyBuffer := bytes.NewBuffer(body)

	req, _ := http.NewRequest("POST", "https://api.mailjet.com/v3/REST/contact", bodyBuffer)
	req.SetBasicAuth(apiKey, apiSecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.StdErr.Println(err)
		return
	}
	defer resp.Body.Close()

	result := struct {
		Data []struct {
			ID int
		}
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logger.StdErr.Println(err)
		return
	}

	contactId := result.Data[0].ID

	// Update contact metadata
	// PUT https://api.mailjet.com/v3/REST/contactdata/$contactId
	// { "Data": [{"Name": "firstname", "Value":"first name!"}] }
	body, _ = json.Marshal(bson.M{
		"Data": bson.A{
			bson.M{
				"Name":  "firstname",
				"Value": firstName,
			},
			bson.M{
				"Name":  "lastname",
				"Value": lastName,
			},
			bson.M{
				"Name":  "picture",
				"Value": picture,
			},
		},
	})
	bodyBuffer = bytes.NewBuffer(body)

	req, _ = http.NewRequest("PUT", fmt.Sprintf("https://api.mailjet.com/v3/REST/contactdata/%d", contactId), bodyBuffer)
	req.SetBasicAuth(apiKey, apiSecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		logger.StdErr.Println(err)
		return
	}
	defer resp.Body.Close()

	// Add contact to "schej users" contact list
	// POST https://api.mailjet.com/v3/REST/contact/$contactId/managecontactslists
	// '{ "ContactsLists": [{"Action": "addforce", "ListID": "10219365"}] }'
	body, _ = json.Marshal(bson.M{
		"ContactsLists": bson.A{
			bson.M{
				"Action": "addforce",
				"ListID": listId,
			},
		},
	})
	bodyBuffer = bytes.NewBuffer(body)

	req, _ = http.NewRequest("POST", fmt.Sprintf("https://api.mailjet.com/v3/REST/contact/%d/managecontactslists", contactId), bodyBuffer)
	req.SetBasicAuth(apiKey, apiSecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		logger.StdErr.Println(err)
		return
	}
	defer resp.Body.Close()
}
