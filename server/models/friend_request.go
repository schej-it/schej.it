package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Representation of a friend request in the MongoDB database
type FriendRequest struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	From      primitive.ObjectID `json:"from" bson:"from,omitempty"`
	FromUser  *User              `json:"fromUser" bson:",omitempty"`
	To        primitive.ObjectID `json:"to" bson:"to,omitempty"`
	ToUser    *User              `json:"toUser" bson:",omitempty"`
	CreatedAt primitive.DateTime `json:"createdAt" bson:"createdAt,omitempty"`
}
