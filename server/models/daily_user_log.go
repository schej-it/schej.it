package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DailyUserLog struct {
	Id      primitive.ObjectID   `json:"_id" bson:"_id,omitempty"`
	Date    primitive.DateTime   `json:"date" bson:"date,omitempty"`
	UserIds []primitive.ObjectID `json:"-" bson:"userIds"`
	Users   []User               `json:"users" bson:",omitempty"`
}
