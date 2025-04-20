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
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27018"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	// Get collections
	eventsCollection := client.Database("schej-it").Collection("events")
	// eventResponsesCollection := client.Database("schej-it").Collection("eventResponses")

	// Get all events
	latestID, err := primitive.ObjectIDFromHex("6804a0d136c40b06cf27aca9")
	if err != nil {
		log.Fatal(err)
	}

	// earliestID, err := primitive.ObjectIDFromHex("67f7e7e39ddc87da36eec9e3")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Create a pipeline to get event IDs and their response counts
	pipeline := []bson.M{
		{
			"$match": bson.M{
				// "_id": bson.M{"$lt": earliestID},
				"_id": bson.M{"$gt": latestID},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "eventResponses",
				"localField":   "_id",
				"foreignField": "eventId",
				"as":           "responses",
			},
		},
		{
			"$project": bson.M{
				"_id":          1,
				"numResponses": bson.M{"$size": "$responses"},
			},
		},
	}

	cursor, err := eventsCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	// Process in batches of 1000
	batchSize := 1000
	var updates []mongo.WriteModel
	var processedCount int

	for cursor.Next(context.Background()) {
		var result struct {
			ID           primitive.ObjectID `bson:"_id"`
			NumResponses int                `bson:"numResponses"`
		}
		if err := cursor.Decode(&result); err != nil {
			log.Printf("Error decoding document: %v", err)
			continue
		}

		update := mongo.NewUpdateOneModel().
			SetFilter(bson.M{"_id": result.ID}).
			SetUpdate(bson.M{"$set": bson.M{"numResponses": result.NumResponses}})
		updates = append(updates, update)

		if len(updates) >= batchSize {
			_, err := eventsCollection.BulkWrite(context.Background(), updates)
			if err != nil {
				log.Printf("Error executing bulk update: %v", err)
			} else {
				processedCount += len(updates)
				fmt.Printf("Processed %d events\n", processedCount)
			}
			updates = updates[:0] // Clear the updates slice
		}
	}

	// Process any remaining updates
	if len(updates) > 0 {
		_, err := eventsCollection.BulkWrite(context.Background(), updates)
		if err != nil {
			log.Printf("Error executing final bulk update: %v", err)
		} else {
			processedCount += len(updates)
			fmt.Printf("Processed final batch of %d events\n", len(updates))
		}
	}

	fmt.Printf("Migration completed! Total events processed: %d\n", processedCount)
}
