package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Representation of an Event in the mongoDB database
type Event struct {
	Id      primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	OwnerId primitive.ObjectID `json:"ownerId" bson:"ownerId,omitempty"`
	Name    string             `json:"name" bson:"name,omitempty"`

	Duration *float32             `json:"duration" bson:"duration,omitempty"`
	Dates    []primitive.DateTime `json:"dates" bson:"dates,omitempty"`

	// Availability responses
	Responses map[string]*Response `json:"responses" bson:"responses"`
}

// A response object containing an array of times that the given user is available
type Response struct {
	Name         string             `json:"name" bson:",omitempty"`
	UserId       primitive.ObjectID `json:"userId" bson:",omitempty"`
	User         *UserProfile       `json:"user" bson:",omitempty"`
	Availability []string           `json:"availability" bson:"availability"`
}
