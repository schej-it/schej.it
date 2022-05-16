package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Event struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" binding:"required"`
	Name      string             `bson:"name,omitempty" binding:"required"`
	StartDate string             `bson:"start_date,omitempty" binding:"required"`
	EndDate   string             `bson:"end_date,omitempty" binding:"required"`
	StartTime int                `bson:"start_time,omitempty" binding:"required"`
	EndTime   int                `bson:"end_time,omitempty" binding:"required"`
	Responses []Response         `bson:"responses,omitempty" binding:"required"`
}

type Response struct {
	Name  string   `bson:"name,omitempty" binding:"required"`
	Times []string `bson:"times,omitempty" binding:"required"`
}
