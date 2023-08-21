package commands

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"schej.it/server/db"
	"schej.it/server/logger"
)

var numUsers Command = Command{
	Name:        "/num_users",
	Description: "Returns the number of signed up users",
	Execute: func(args []string, webhookUrl string) {
		var results []bson.M
		cursor, err := db.UsersCollection.Aggregate(context.Background(), []bson.M{
			{"$match": bson.M{}},
			{"$group": bson.M{"_id": nil, "n": bson.M{"$sum": 1}}},
		})
		if err != nil {
			logger.StdErr.Panicln(err)
		}
		if err := cursor.All(context.Background(), &results); err != nil {
			logger.StdErr.Panicln(err)
		}
		n := results[0]["n"]

		response := Response{
			ResponseType: "in_channel",
			Text:         fmt.Sprintf("Number of currently signed up users: %v", n),
		}
		SendRawMessage(&response, webhookUrl)
	},
}
