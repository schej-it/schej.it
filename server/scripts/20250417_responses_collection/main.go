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
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27018"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	// Get collections
	eventsCollection := client.Database("schej-it").Collection("events")
	eventResponsesCollection := client.Database("schej-it").Collection("eventResponses")
	attendeesCollection := client.Database("schej-it").Collection("attendees")

	// Get all events
	lastProcessedID, err := primitive.ObjectIDFromHex("6804a03836c40b06cf27aca4")
	if err != nil {
		log.Fatal(err)
	}
	cursor, err := eventsCollection.Find(
		context.Background(),
		bson.M{
			"_id": bson.M{"$gt": lastProcessedID},
		},
		options.Find().SetSort(bson.M{"_id": 1}),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	// Iterate through events
	eventResponses := make([]interface{}, 0)
	eventIds := make([]primitive.ObjectID, 0)
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
					EventId:  eventId,
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
			eventResponses = append(eventResponses, eventResponse)

			// Insert into eventResponses collection
			// _, err := eventResponsesCollection.InsertOne(context.Background(), eventResponse)
			// if err != nil {
			// 	log.Printf("Error inserting response for event %s, user %s: %v", eventId, userIdString, err)
			// 	continue
			// }

			// fmt.Printf("Migrated response for event %s, user %s\n", eventId, userIdString)
		}

		// Insert into eventResponses collection
		eventIds = append(eventIds, eventId)
		if len(eventResponses) > 200 {
			_, err := eventResponsesCollection.InsertMany(context.Background(), eventResponses)
			if err != nil {
				log.Printf("Error inserting responses for event %s: %v", eventId, err)
			}
			fmt.Printf("Inserted %d responses for eventIds %v\n", len(eventResponses), eventIds)
			eventResponses = make([]interface{}, 0)
			eventIds = make([]primitive.ObjectID, 0)
		}
	}

	if len(eventResponses) > 0 {
		_, err := eventResponsesCollection.InsertMany(context.Background(), eventResponses)
		if err != nil {
			log.Printf("Error inserting responses for eventIds %v: %v", eventIds, err)
		}
		fmt.Printf("Inserted %d responses for eventIds %v\n", len(eventResponses), eventIds)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Migration completed successfully!")
}
