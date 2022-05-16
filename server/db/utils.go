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

	/*
		cursor, err := EventsCollection.Aggregate(context.Background(), []bson.M{
			{"$match": bson.M{"_id": utils.StringToObjectID(eventId)}},
			{"$objectToArray": "$responses"},
			{"$unwind": "$responses"},
			{"$lookup": bson.M{
				"from":         "users",
				"localField":   "responses.userId",
				"foreignField": "_id",
				"as":           "responses.user",
			}},
		})
		if err != nil {
			panic(err)
		}
		var events []models.Event
		if err := cursor.All(context.Background(), &events); err != nil {
			panic(err)
		}

		if len(events) == 0 {
			// Event does not exist!
			return nil
		}

		return &events[0]
	*/
}
