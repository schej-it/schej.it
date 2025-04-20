package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Attendee struct {
	Id      primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	EventId primitive.ObjectID `json:"eventId" bson:"eventId,omitempty"`

	Email    string `json:"email" bson:"email,omitempty"`
	Declined *bool  `json:"declined" bson:"declined,omitempty"`
}
