package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Event struct {
	Id        primitive.ObjectID  `json:"_id" bson:"_id,omitempty" binding:"required"`
	OwnerId   primitive.ObjectID  `json:"ownerId" bson:"ownerId" binding:"required"`
	Name      string              `json:"name" bson:"name" binding:"required"`
	StartDate primitive.DateTime  `json:"startDate" bson:"startDate" binding:"required"`
	EndDate   primitive.DateTime  `json:"endDate" bson:"endDate" binding:"required"`
	StartTime int                 `json:"startTime" bson:"startTime" binding:"required"`
	EndTime   int                 `json:"endTime" bson:"endTime" binding:"required"`
	Responses map[string]Response `json:"responses" bson:"responses" binding:"required"`
}

type Response struct {
	UserId       primitive.ObjectID `json:"userId" bson:"userId" binding:"required"`
	Availability []string           `json:"availability" bson:"availability" binding:"required"`
}
