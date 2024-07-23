package calendar

import (
	"time"

	"schej.it/server/models"
)

type CalendarProvider interface {
	GetEmail() string
	GetCalendarList() (map[string]models.SubCalendar, error)
	GetCalendarEvents(calendarId string, timeMin time.Time, timeMax time.Time) ([]models.CalendarEvent, error)
}

func GetCalendarProvider(calendarAccount models.CalendarAccount) CalendarProvider {
	switch calendarAccount.CalendarType {
	case models.GoogleCalendarType:
		return &GoogleCalendar{
			GoogleCalendarDetails: *calendarAccount.GoogleCalendarDetails,
		}
	case models.AppleCalendarType:
		return &AppleCalendar{
			AppleCalendarDetails: *calendarAccount.AppleCalendarDetails,
		}
	}
	return nil
}
