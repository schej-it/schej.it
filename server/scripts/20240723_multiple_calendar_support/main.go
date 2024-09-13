package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"schej.it/server/db"
	"schej.it/server/models"
	"schej.it/server/utils"
)

type OldUser struct {
	TimezoneOffset int `json:"timezoneOffset" bson:"timezoneOffset"`

	// Profile info
	Id        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Email     string             `json:"email" bson:"email,omitempty"`
	FirstName string             `json:"firstName" bson:"firstName,omitempty"`
	LastName  string             `json:"lastName" bson:"lastName,omitempty"`
	Picture   string             `json:"picture" bson:"picture,omitempty"`

	HasCustomName *bool `json:"hasCustomName" bson:"hasCustomName,omitempty"`

	// Changed from CalendarAccount to OldCalendarAccount
	CalendarAccounts map[string]OldCalendarAccount `json:"calendarAccounts" bson:"calendarAccounts,omitempty"`

	TokenOrigin models.TokenOriginType `json:"-" bson:"tokenOrigin,omitempty"`

	CalendarOptions *models.CalendarOptions `json:"calendarOptions" bson:"calendarOptions,omitempty"`
}

type OldCalendarAccount struct {
	Email        string                         `bson:"email,omitempty"`
	Picture      string                         `bson:"picture,omitempty"`
	Enabled      *bool                          `bson:"enabled,omitempty"`
	SubCalendars *map[string]models.SubCalendar `bson:"subCalendars,omitempty"`

	AccessToken           string             `bson:"accessToken,omitempty"`
	AccessTokenExpireDate primitive.DateTime `bson:"accessTokenExpireDate,omitempty"`
	RefreshToken          string             `bson:"refreshToken,omitempty"`
}

func main() {
	closeConnection := db.Init()
	defer closeConnection()

	cursor, err := db.UsersCollection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var oldUser OldUser
		if err := cursor.Decode(&oldUser); err != nil {
			log.Printf("Error decoding user: %v", err)
			continue
		}

		newCalendarAccounts := make(map[string]models.CalendarAccount)
		for _, calendarAccount := range oldUser.CalendarAccounts {
			newCalendarAccount := models.CalendarAccount{
				CalendarType: models.GoogleCalendarType,
				OAuth2CalendarAuth: &models.OAuth2CalendarAuth{
					AccessToken:           calendarAccount.AccessToken,
					AccessTokenExpireDate: calendarAccount.AccessTokenExpireDate,
					RefreshToken:          calendarAccount.RefreshToken,
				},

				Email:        calendarAccount.Email,
				Picture:      calendarAccount.Picture,
				Enabled:      calendarAccount.Enabled,
				SubCalendars: calendarAccount.SubCalendars,
			}
			calendarAccountKey := utils.GetCalendarAccountKey(calendarAccount.Email, models.GoogleCalendarType)
			newCalendarAccounts[calendarAccountKey] = newCalendarAccount
		}

		update := bson.M{
			"$set": bson.M{
				"calendarAccounts": newCalendarAccounts,
			},
		}

		_, err := db.UsersCollection.UpdateByID(context.Background(), oldUser.Id, update)
		if err != nil {
			log.Printf("Error updating user %s: %v", oldUser.Email, err)
		} else {
			log.Printf("Successfully updated user: %s", oldUser.Email)
		}
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	log.Println("Migration completed successfully")
}
