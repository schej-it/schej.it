package utils

import (
	"go.mongodb.org/mongo-driver/bson"
	"schej.it/server/models"
)

func UpdateEventResponseAggregation(userIdString string, response models.Response) bson.M {
	res := bson.M{"$set": bson.M{
		"responses": bson.M{
			"$setField": bson.M{
				"field": bson.M{"$literal": userIdString},                        // Use $setField because userIdString could have periods
				"input": bson.M{"$ifNull": bson.A{"$$ROOT.responses", bson.M{}}}, // If responses map is null, create it
				"value": bson.M{"$literal": response},
			},
		},
	}}
	return res
}

func DeleteEventResponseAggregation(userIdString string) bson.M {
	res := bson.M{"$set": bson.M{
		"responses": bson.M{
			"$setField": bson.M{
				"field": bson.M{"$literal": userIdString},                        // Use $setField because userIdString could have periods
				"input": bson.M{"$ifNull": bson.A{"$$ROOT.responses", bson.M{}}}, // If responses map is null, create it
				"value": "$$REMOVE",
			},
		},
	}}
	return res
}
