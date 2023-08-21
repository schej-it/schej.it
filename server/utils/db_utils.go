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
