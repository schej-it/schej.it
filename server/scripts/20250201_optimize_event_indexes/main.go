package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"schej.it/server/db"
)

func main() {
	// Initialize database connection
	disconnect := db.Init()
	defer disconnect()

	// Drop old indexes if they exist
	_, err := db.EventsCollection.Indexes().DropOne(
		context.Background(),
		"responses_userId_1",
	)
	if err != nil {
		fmt.Printf("Warning: Failed to drop old responses index: %v\n", err)
	}

	_, err = db.EventsCollection.Indexes().DropOne(
		context.Background(),
		"attendees_email_1",
	)
	if err != nil {
		fmt.Printf("Warning: Failed to drop old attendees index: %v\n", err)
	}

	// Create compound index on responses.userId and _id for better query performance with sorting
	_, err = db.EventsCollection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.D{
				{Key: "responses.userId", Value: 1},
				{Key: "_id", Value: -1},
			},
			Options: options.Index().
				SetName("responses_userId_id_1"),
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Created compound index on responses.userId and _id")

	// Create compound index on attendees.email, declined status, and _id
	_, err = db.EventsCollection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.D{
				{Key: "attendees.email", Value: 1},
				{Key: "attendees.declined", Value: 1},
				{Key: "_id", Value: -1},
			},
			Options: options.Index().
				SetName("attendees_email_declined_id_1"),
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Created compound index on attendees.email, declined status, and _id")

	os.Exit(0)
}
