package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Simplified representation of a Google Calendar event from the calendar api
type CalendarEvent struct {
	Summary   string             `json:"summary" bson:"summary,omitempty"`
	StartDate primitive.DateTime `json:"startDate" bson:"startDate,omitempty"`
	EndDate   primitive.DateTime `json:"endDate" bson:"endDate,omitempty"`
}
