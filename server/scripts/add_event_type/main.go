package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"schej.it/server/db"
	"schej.it/server/models"
)

func main() {
	closeConnection := db.Init()
	defer closeConnection()

	_, err := db.EventsCollection.UpdateMany(context.Background(), bson.M{"type": nil}, bson.M{
		"$set": bson.M{"type": models.SPECIFIC_DATES},
	})
	if err != nil {
		log.Fatal(err)
	}

}
