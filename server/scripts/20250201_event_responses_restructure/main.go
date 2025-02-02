package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"schej.it/server/db"
	"schej.it/server/models"
)

// OldEvent represents the event structure before the migration
type OldEvent struct {
	Id        primitive.ObjectID          `bson:"_id,omitempty"`
	Responses map[string]*models.Response `bson:"responses"`
	OwnerId   primitive.ObjectID          `bson:"ownerId,omitempty"`
}

// NewEvent represents the event structure after the migration
type NewEvent struct {
	Id            primitive.ObjectID     `bson:"_id,omitempty"`
	ResponsesList []models.EventResponse `bson:"responses"`
	OwnerId       primitive.ObjectID     `bson:"ownerId,omitempty"`
}

func main() {
	// Initialize database connection
	disconnect := db.Init()
	defer disconnect()

	// Get all events
	cursor, err := db.EventsCollection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var oldEvents []OldEvent
	if err = cursor.All(context.Background(), &oldEvents); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d events to migrate\n", len(oldEvents))

	// Process events in batches
	batchSize := 100
	totalUpdated := 0

	for i := 0; i < len(oldEvents); i += batchSize {
		end := i + batchSize
		if end > len(oldEvents) {
			end = len(oldEvents)
		}

		var operations []mongo.WriteModel
		for _, oldEvent := range oldEvents[i:end] {
			// Skip if event has no responses
			if oldEvent.Responses == nil {
				continue
			}

			// Convert map to array format
			var responsesList []models.EventResponse
			for userIdHex, response := range oldEvent.Responses {
				if err != nil {
					fmt.Printf("Warning: Invalid userId hex %s, skipping\n", userIdHex)
					continue
				}
				responsesList = append(responsesList, models.EventResponse{
					UserId:   userIdHex,
					Response: response,
				})
			}

			// Create update operation
			update := mongo.NewUpdateOneModel()
			update.SetFilter(bson.M{"_id": oldEvent.Id})
			update.SetUpdate(bson.M{
				"$set": bson.M{
					"responses": responsesList,
				},
				"$unset": bson.M{
					"responsesMap": "",
				},
			})
			operations = append(operations, update)
		}

		// Execute batch update
		if len(operations) > 0 {
			result, err := db.EventsCollection.BulkWrite(context.Background(), operations)
			if err != nil {
				log.Fatal(err)
			}
			totalUpdated += int(result.ModifiedCount)
			fmt.Printf("Updated %d events in batch\n", result.ModifiedCount)
		}
	}

	fmt.Printf("Migration complete. Updated %d events total\n", totalUpdated)

	// Create index on responses.userId
	_, err = db.EventsCollection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.D{{Key: "responses.userId", Value: 1}},
			Options: options.Index().
				SetName("responses_userId_1"),
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Created index on responses.userId")

	// Create index on attendees.email
	_, err = db.EventsCollection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.D{{Key: "attendees.email", Value: 1}},
			Options: options.Index().
				SetName("attendees_email_1"),
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Created index on attendees.email")

	os.Exit(0)
}
