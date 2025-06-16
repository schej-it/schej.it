package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Folder struct {
	Id     primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserId primitive.ObjectID `json:"userId" bson:"userId"`

	Name      string  `json:"name,omitempty" bson:"name,omitempty"`
	Color     *string `json:"color,omitempty" bson:"color,omitempty"`
	IsDeleted *bool   `json:"isDeleted,omitempty" bson:"isDeleted,omitempty"`

	EventIds []primitive.ObjectID `json:"eventIds" bson:"-"`
}
