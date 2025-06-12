package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Folder struct {
	Id       primitive.ObjectID  `json:"_id" bson:"_id,omitempty"`
	UserId   primitive.ObjectID  `json:"userId" bson:"userId"`
	ParentId *primitive.ObjectID `json:"parentId,omitempty" bson:"parentId,omitempty"`

	Name      string `json:"name,omitempty" bson:"name,omitempty"`
	IsDeleted *bool  `json:"isDeleted,omitempty" bson:"isDeleted,omitempty"`

	Folders []Folder `json:"folders" bson:"-"`
	Events  []Event  `json:"events" bson:"-"`
}
