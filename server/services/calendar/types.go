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
