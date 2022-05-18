package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Simplified representation of a Google Calendar event from the calendar api
type CalendarEvent struct {
	Summary   string             `json:"summary" bson:"summary" binding:"required"`
	StartDate primitive.DateTime `json:"startDate" bson:"startDate" binding:"required"`
	EndDate   primitive.DateTime `json:"endDate" bson:"endDate" binding:"required"`
}
