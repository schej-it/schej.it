package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" binding:"required"`
	Email     string             `bson:"email,omitempty" binding:"required"`
	FirstName string             `bson:"first_name,omitempty" binding:"required"`
	LastName  string             `bson:"last_name,omitempty" binding:"required"`
	Picture   string             `bson:"picture,omitempty" binding:"required"`
}
