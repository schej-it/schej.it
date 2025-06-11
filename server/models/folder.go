package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Folder struct {
	Id       primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserId   primitive.ObjectID `json:"userId" bson:"userId"`
	ParentId *string            `json:"parentId" bson:"parentId,omitempty"`

	Name string `json:"name" bson:"name,omitempty"`

	Folders []primitive.ObjectID `json:"folders" bson:"folders,omitempty"`
	Events  []primitive.ObjectID `json:"events" bson:"events,omitempty"`

	IsDeleted *bool `json:"isDeleted" bson:"isDeleted,omitempty"`
}
