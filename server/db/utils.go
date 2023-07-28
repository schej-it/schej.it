package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"schej.it/server/logger"
	"schej.it/server/models"
	"schej.it/server/services/auth"
	"schej.it/server/utils"
)

func GetUserById(userId string) *models.User {
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		// userId is malformatted
		return nil
	}
	result := UsersCollection.FindOne(context.Background(), bson.M{
		"_id": objectId,
	})
	if result.Err() == mongo.ErrNoDocuments {
		// User does not exist!
		return nil
	}

	// Decode result
	var user models.User
	if err := result.Decode(&user); err != nil {
		logger.StdErr.Panicln(err)
	}

	return &user
}

func GetEventById(eventId string) *models.Event {
	objectId, err := primitive.ObjectIDFromHex(eventId)
	if err != nil {
		// eventId is malformatted
		return nil
	}
	result := EventsCollection.FindOne(context.Background(), bson.M{
		"_id": objectId,
	})
	if result.Err() == mongo.ErrNoDocuments {
		// Event does not exist!
		return nil
	}

	// Decode result
	var event models.Event
	if err := result.Decode(&event); err != nil {
		logger.StdErr.Panicln(err)
	}

	return &event
}

func GetFriendRequestById(friendRequestId string) *models.FriendRequest {
	objectId, err := primitive.ObjectIDFromHex(friendRequestId)
	if err != nil {
		// friendRequestId is malformatted
		return nil
	}
	result := FriendRequestsCollection.FindOne(context.Background(), bson.M{
		"_id": objectId,
	})
	if result.Err() == mongo.ErrNoDocuments {
		// Friend request does not exist!
		return nil
	}

	// Decode result
	var friendRequest models.FriendRequest
	if err := result.Decode(&friendRequest); err != nil {
		logger.StdErr.Panicln(err)
	}

	return &friendRequest
}

func DeleteFriendRequestById(friendRequestId string) {
	objectId, err := primitive.ObjectIDFromHex(friendRequestId)
	if err != nil {
		// friendRequestId is malformatted
		logger.StdErr.Panicln(err)
	}
	_, err = FriendRequestsCollection.DeleteOne(context.Background(), bson.M{
		"_id": objectId,
	})
	if err != nil {
		logger.StdErr.Panicln(err)
	}
}

// If access token has expired, get a new token for the primary account as well as all other calendar accounts, update the user object, and save it to the database
func RefreshUserTokenIfNecessary(u *models.User) {
	numAccountsToUpdate := 0
	refreshTokenChan := make(chan auth.RefreshAccessTokenData)

	// Refresh primary account access token if necessary
	if time.Now().After(u.AccessTokenExpireDate.Time()) && len(u.RefreshToken) > 0 {
		go auth.RefreshAccessTokenAsync(u.RefreshToken, refreshTokenChan, -1, true)
		numAccountsToUpdate++
	}

	// Refresh other calendar account access tokens if necessary
	for i, account := range u.CalendarAccounts {
		if time.Now().After(account.AccessTokenExpireDate.Time()) && len(account.RefreshToken) > 0 {
			go auth.RefreshAccessTokenAsync(account.RefreshToken, refreshTokenChan, i, false)
			numAccountsToUpdate++
		}
	}

	// Update access tokens as responses are received
	for i := 0; i < numAccountsToUpdate; i++ {
		res := <-refreshTokenChan

		accessTokenExpireDate := utils.GetAccessTokenExpireDate(res.TokenResponse.ExpiresIn)

		if res.IsPrimaryAccount {
			u.AccessToken = res.TokenResponse.AccessToken
			u.AccessTokenExpireDate = primitive.NewDateTimeFromTime(accessTokenExpireDate)
		} else {
			u.CalendarAccounts[res.CalendarAccountIndex].AccessToken = res.TokenResponse.AccessToken
			u.CalendarAccounts[res.CalendarAccountIndex].AccessTokenExpireDate = primitive.NewDateTimeFromTime(accessTokenExpireDate)
		}
	}

	// Update user object if accounts were updated
	if numAccountsToUpdate > 0 {
		UsersCollection.FindOneAndUpdate(
			context.Background(),
			bson.M{"_id": u.Id},
			bson.M{"$set": u},
		)
	}
}

/*
Finds a daily user log, localized to the user's month/day/year, by:
Checking if the given date (in server time) and timezone offset (in client time) match the day of a daily user log (in UTC time)
If none exist, then create a new daily log

Note: we find the log localized to the user's month/day/year in order to track if the user has signed in on different days in their
own timezone, rather than the server's timezone. For example, if a user signed in at 11pm on Monday, then signed in at 8am on Tuesday,
it could theoretically count as the same day if we were to use server time
*/
func GetDailyUserLogByDate(date time.Time, timezoneOffset int) *models.DailyUserLog {
	timezoneOffsetDuration, _ := time.ParseDuration(fmt.Sprintf("%dm", timezoneOffset))
	adjustedDate := date.Add(timezoneOffsetDuration)
	startDate := utils.GetDateAtTime(adjustedDate, "00:00:00")
	endDate := utils.GetDateAtTime(adjustedDate, "23:59:59")

	// Find a log for the current date
	result := DailyUserLogCollection.FindOne(context.Background(), bson.M{
		"date": bson.M{
			"$gte": primitive.NewDateTimeFromTime(startDate),
			"$lte": primitive.NewDateTimeFromTime(endDate),
		},
	})

	var log models.DailyUserLog

	// Create a new log if it doesn't exist already
	if result.Err() == mongo.ErrNoDocuments {
		log = models.DailyUserLog{
			Date: primitive.NewDateTimeFromTime(startDate),
		}
		result, err := DailyUserLogCollection.InsertOne(context.Background(), log)
		if err != nil {
			logger.StdErr.Panicln(err)
		}
		log.Id = result.InsertedID.(primitive.ObjectID)
	} else {
		// Parse daily user log object
		if err := result.Decode(&log); err != nil {
			logger.StdErr.Panicln(err)
		}
	}

	return &log
}

func UpdateDailyUserLog(user *models.User) {
	log := GetDailyUserLogByDate(time.Now(), user.TimezoneOffset)
	for _, id := range log.UserIds {
		if id == user.Id {
			return
		}
	}

	log.UserIds = append(log.UserIds, user.Id)
	_, err := DailyUserLogCollection.UpdateByID(context.Background(), log.Id, bson.M{"$set": log})
	if err != nil {
		logger.StdErr.Panicln(err)
	}
}
