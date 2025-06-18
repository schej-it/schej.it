package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"schej.it/server/logger"
)

func CountDistinctMonthlyActiveEventCreators(date time.Time) (int64, error) {
	thirtyDaysAgo := date.AddDate(0, 0, -30)
	// Generate a minimal ObjectID for the timestamp 30 days ago
	minObjectId := primitive.NewObjectIDFromTimestamp(thirtyDaysAgo)
	maxObjectId := primitive.NewObjectIDFromTimestamp(date)

	filter := bson.M{
		"_id": bson.M{
			"$gte": minObjectId,
			"$lte": maxObjectId,
		},
		// Ensure creatorPosthogId exists and is not null/empty if needed
		"creatorPosthogId": bson.M{"$exists": true, "$ne": ""},
	}

	distinctValues, err := EventsCollection.Distinct(context.Background(), "creatorPosthogId", filter)
	if err != nil {
		logger.StdErr.Printf("Error counting distinct monthly active creators: %v\n", err)
		return 0, err
	}

	return int64(len(distinctValues)), nil
}

func CountDistinctMonthlyActiveEventCreatorsWithMoreThanXEvents(date time.Time, x int) (int64, error) {
	thirtyDaysAgo := date.AddDate(0, 0, -30)
	// Generate a minimal ObjectID for the timestamp 30 days ago
	minObjectId := primitive.NewObjectIDFromTimestamp(thirtyDaysAgo)
	maxObjectId := primitive.NewObjectIDFromTimestamp(date)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"_id": bson.M{
				"$gte": minObjectId,
				"$lte": maxObjectId,
			},
			"creatorPosthogId": bson.M{"$exists": true, "$ne": ""},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":   "$creatorPosthogId",
			"count": bson.M{"$sum": 1},
		}}},
		{{Key: "$match", Value: bson.M{
			"count": bson.M{"$gte": x},
		}}},
		{{Key: "$count", Value: "creatorCount"}},
	}

	cursor, err := EventsCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		logger.StdErr.Printf("Error aggregating active creators with more than %d events: %v\n", x, err)
		return 0, err
	}
	defer cursor.Close(context.Background())

	var results []bson.M
	if err := cursor.All(context.Background(), &results); err != nil {
		logger.StdErr.Printf("Error decoding aggregation result: %v\n", err)
		return 0, err
	}

	if len(results) == 0 {
		return 0, nil
	}

	if count, ok := results[0]["creatorCount"].(int32); ok {
		return int64(count), nil
	}
	if count, ok := results[0]["creatorCount"].(int64); ok {
		return count, nil
	}

	return 0, fmt.Errorf("could not convert creatorCount to int64")
}
