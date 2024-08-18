package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"schej.it/server/db"
)

func main() {
	closeConnection := db.Init()
	defer closeConnection()

	_, err := db.EventsCollection.UpdateMany(context.Background(), bson.M{}, bson.M{
		"$rename": bson.M{"blindavailabilityenabled": "blindAvailabilityEnabled"},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Done!!")
}
