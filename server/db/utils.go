package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"schej.it/server/models"
	"schej.it/server/utils"
)

func GetUserById(userId string) *models.User {
	result := UsersCollection.FindOne(context.Background(), bson.M{
		"_id": utils.StringToObjectID(userId),
	})
	if result.Err() == mongo.ErrNoDocuments {
		// User does not exist!
		return nil
	}

	// Set auth user request variable
	var user models.User
	if err := result.Decode(&user); err != nil {
		panic(err)
	}

	return &user
}

func GetEventById(eventId string) *models.Event {
	result := EventsCollection.FindOne(context.Background(), bson.M{
		"_id": utils.StringToObjectID(eventId),
	})
	if result.Err() == mongo.ErrNoDocuments {
		// Event does not exist!
		return nil
	}

	// Set auth user request variable
	var event models.Event
	if err := result.Decode(&event); err != nil {
		panic(err)
	}

	return &event
}
