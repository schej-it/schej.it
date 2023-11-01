package utils

import (
	"go.mongodb.org/mongo-driver/bson"
	"schej.it/server/models"
)

// Returns an aggregation expression inserting/updating a calendar account for the user
func InsertCalendarAccountAggregation(calendarAccount models.CalendarAccount) bson.M {
	return bson.M{"$set": bson.M{
		"calendarAccounts": bson.M{
			"$setField": bson.M{
				"field": calendarAccount.Email,                                          // Use $setField because email will have periods
				"input": bson.M{"$ifNull": bson.A{"$$ROOT.calendarAccounts", bson.M{}}}, // If calendarAccounts map is null, create it
				"value": calendarAccount,
			},
		},
	}}
}

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
