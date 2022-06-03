package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Representation of an Event in the mongoDB database
type Event struct {
	Id        primitive.ObjectID  `json:"_id" bson:"_id,omitempty"`
	OwnerId   primitive.ObjectID  `json:"ownerId" bson:"ownerId,omitempty"`
	Name      string              `json:"name" bson:"name,omitempty"`
	StartDate primitive.DateTime  `json:"startDate" bson:"startDate,omitempty"`
	EndDate   primitive.DateTime  `json:"endDate" bson:"endDate,omitempty"`
	Responses map[string]Response `json:"responses" bson:"responses,omitempty"`
}

// A response object containing an array of times that the given user is available
type Response struct {
	UserId       primitive.ObjectID `json:"userId" bson:"userId,omitempty"`
	User         *UserProfile       `json:"user" bson:",omitempty"`
	Availability []string           `json:"availability" bson:"availability,omitempty"`
}
