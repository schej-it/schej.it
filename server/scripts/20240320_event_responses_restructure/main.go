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

func main() {
	// Initialize database connection
	disconnect := db.Init()
	defer disconnect()

	// Get all events
	cursor, err := db.EventsCollection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var events []models.Event
	if err = cursor.All(context.Background(), &events); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d events to migrate\n", len(events))

	// Process events in batches
	batchSize := 100
	totalUpdated := 0

	for i := 0; i < len(events); i += batchSize {
		end := i + batchSize
		if end > len(events) {
			end = len(events)
		}

		var operations []mongo.WriteModel
		for _, event := range events[i:end] {
			// Skip if event has no responses
			if event.ResponsesMap == nil && len(event.ResponsesList) == 0 {
				continue
			}

			// Convert map to array format if needed
			var responsesList []models.EventResponse
			if len(event.ResponsesList) == 0 && event.ResponsesMap != nil {
				for userIdHex, response := range event.ResponsesMap {
					userId, err := primitive.ObjectIDFromHex(userIdHex)
					if err != nil {
						fmt.Printf("Warning: Invalid userId hex %s, skipping\n", userIdHex)
						continue
					}
					responsesList = append(responsesList, models.EventResponse{
						UserId:   userId,
						Response: response,
					})
				}
			}

			// Create update operation
			update := mongo.NewUpdateOneModel()
			update.SetFilter(bson.M{"_id": event.Id})
			update.SetUpdate(bson.M{
				"$set": bson.M{
					"responses":    responsesList,
					"responsesMap": event.ResponsesMap,
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

	os.Exit(0)
}
