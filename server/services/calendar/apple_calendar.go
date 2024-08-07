package calendar

import (
	"time"

	"schej.it/server/models"
)

type AppleCalendar struct {
	models.AppleCalendarAuth
}

func (calendar *AppleCalendar) GetCalendarList() (map[string]models.SubCalendar, error) {
	return nil, nil
}

func (calendar *AppleCalendar) GetCalendarEvents(calendarId string, timeMin time.Time, timeMax time.Time) ([]models.CalendarEvent, error) {
	return nil, nil
}
