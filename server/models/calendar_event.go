package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Calendar contains information about a user's calendar
type Calendar struct {
	Id       string `json:"id" bson:"id,omitempty"`
	Summary  string `json:"summary" bson:"summary,omitempty"`
	Selected bool   `json:"selected" bson:"selected,omitempty"`
}

// Simplified representation of a Google Calendar event from the calendar api
type CalendarEvent struct {
	Id         string             `json:"id" bson:"id,omitempty"`
	CalendarId string             `json:"calendarId" bson:"calendarId,omitempty"`
	Summary    string             `json:"summary" bson:"summary,omitempty"`
	StartDate  primitive.DateTime `json:"startDate" bson:"startDate,omitempty"`
	EndDate    primitive.DateTime `json:"endDate" bson:"endDate,omitempty"`
}
