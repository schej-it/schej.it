package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	// Get collections
	eventsCollection := client.Database("schej-it").Collection("events")
	eventResponsesCollection := client.Database("schej-it").Collection("eventResponses")

	// Get all events
	// lastProcessedID, err := primitive.ObjectIDFromHex("676f772e48f2022f2544011a")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	cursor, err := eventsCollection.Find(context.Background(), bson.M{
		// "_id": bson.M{"$gt": lastProcessedID},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	// Iterate through events
	for cursor.Next(context.Background()) {
		var event bson.M
		if err := cursor.Decode(&event); err != nil {
			log.Fatal(err)
		}

		eventId := event["_id"].(primitive.ObjectID)

		// Update numResponses for the event
		count, err := eventResponsesCollection.CountDocuments(context.Background(), bson.M{"eventId": eventId})
		if err != nil {
			log.Printf("Error counting responses for event %s: %v", eventId, err)
			continue
		}

		_, err = eventsCollection.UpdateOne(
			context.Background(),
			bson.M{"_id": eventId},
			bson.M{"$set": bson.M{"numResponses": count}},
		)
		if err != nil {
			log.Printf("Error updating numResponses for event %s: %v", eventId, err)
			continue
		}
		fmt.Printf("Updated numResponses to %d for event %s\n", count, eventId)
	}

	fmt.Println("Migration completed successfully!")
}
