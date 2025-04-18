package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"schej.it/server/models"
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
	attendeesCollection := client.Database("schej-it").Collection("attendees")

	// Get all events
	cursor, err := eventsCollection.Find(context.Background(), bson.M{})
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
		if attendeesData, ok := event["attendees"]; ok && attendeesData != nil {
			attendees := attendeesData.(bson.A)
			for _, attendee := range attendees {
				attendeeMap := attendee.(bson.M)
				attendeeEmail := attendeeMap["email"].(string)
				attendeeDeclined := attendeeMap["declined"].(bool)

				attendee := models.Attendee{
					Email:    attendeeEmail,
					Declined: &attendeeDeclined,
				}

				_, err := attendeesCollection.InsertOne(context.Background(), attendee)
				if err != nil {
					log.Printf("Error inserting attendee for event %s, user %s: %v", eventId, attendeeEmail, err)
					continue
				}
			}
		}

		if responsesData, ok := event["responses"]; !ok || responsesData == nil {
			continue
		}
		responses := event["responses"].(bson.A)

		// Iterate through responses
		for _, response := range responses {
			responseMap := response.(bson.M)
			userIdString := responseMap["userId"].(string)

			// Convert bson.M to models.Response
			var responseData models.Response
			var responseBytes []byte
			responseBytes, err = bson.Marshal(responseMap["response"])
			if err != nil {
				log.Printf("Error marshaling response for event %s, user %s: %v", eventId, userIdString, err)
				continue
			}
			if err = bson.Unmarshal(responseBytes, &responseData); err != nil {
				log.Printf("Error unmarshaling response for event %s, user %s: %v", eventId, userIdString, err)
				continue
			}

			// Create event response
			eventResponse := models.EventResponse{
				UserId:   userIdString,
				Response: &responseData,
				EventId:  eventId,
			}

			// Insert into eventResponses collection
			_, err := eventResponsesCollection.InsertOne(context.Background(), eventResponse)
			if err != nil {
				log.Printf("Error inserting response for event %s, user %s: %v", eventId, userIdString, err)
				continue
			}

			fmt.Printf("Migrated response for event %s, user %s\n", eventId, userIdString)
		}
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Migration completed successfully!")
}
