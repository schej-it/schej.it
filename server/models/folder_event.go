package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// FolderEvent represents the mapping between a folder and an event.
type FolderEvent struct {
	Id       primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserId   primitive.ObjectID `json:"userId" bson:"userId"`
	FolderId primitive.ObjectID `json:"folderId" bson:"folderId"`
	EventId  primitive.ObjectID `json:"eventId" bson:"eventId"`
}
