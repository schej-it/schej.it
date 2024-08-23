package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"schej.it/server/db"
	"schej.it/server/models"
)

type GoogleCalendarAuth struct {
	AccessToken           string             `json:"accessToken" bson:"accessToken"`
	RefreshToken          string             `json:"refreshToken" bson:"refreshToken"`
	AccessTokenExpireDate primitive.DateTime `json:"accessTokenExpireDate" bson:"accessTokenExpireDate"`
}

type OldUser struct {
	// Profile info
	Id        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Email     string             `json:"email" bson:"email,omitempty"`
	FirstName string             `json:"firstName" bson:"firstName,omitempty"`
	LastName  string             `json:"lastName" bson:"lastName,omitempty"`
	Picture   string             `json:"picture" bson:"picture,omitempty"`

	// CalendarAccounts is a mapping from {`email_CALENDARTYPE` => CalendarAccount} that contains all the
	// additional accounts the user wants to see google calendar events for
	CalendarAccounts map[string]OldCalendarAccount `json:"calendarAccounts" bson:"calendarAccounts,omitempty"`
}

type OldCalendarAccount struct {
	CalendarType       models.CalendarType       `json:"calendarType" bson:"calendarType,omitempty"`
	GoogleCalendarAuth *GoogleCalendarAuth       `json:"googleCalendarAuth" bson:"googleCalendarAuth,omitempty"`
	AppleCalendarAuth  *models.AppleCalendarAuth `json:"appleCalendarAuth" bson:"appleCalendarAuth,omitempty"`

	Email        string                         `json:"email" bson:"email"` // Email is required for all calendar accounts
	Picture      string                         `json:"picture" bson:"picture,omitempty"`
	Enabled      *bool                          `json:"enabled" bson:"enabled,omitempty"`
	SubCalendars *map[string]models.SubCalendar `json:"subCalendars" bson:"subCalendars,omitempty"`
}

func main() {
	closeConnection := db.Init()
	defer closeConnection()

	cursor, err := db.UsersCollection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var users []OldUser
	if err = cursor.All(context.Background(), &users); err != nil {
		log.Fatal(err)
	}

	for _, user := range users {
		fmt.Printf("Processing user: %s\n", user.Email)

		newCalendarAccounts := make(map[string]models.CalendarAccount)

		for i, oldCalendarAccount := range user.CalendarAccounts {
			newCalendarAccount := models.CalendarAccount{
				CalendarType:      oldCalendarAccount.CalendarType,
				Email:             oldCalendarAccount.Email,
				Picture:           oldCalendarAccount.Picture,
				Enabled:           oldCalendarAccount.Enabled,
				SubCalendars:      oldCalendarAccount.SubCalendars,
				AppleCalendarAuth: oldCalendarAccount.AppleCalendarAuth,
			}

			if oldCalendarAccount.GoogleCalendarAuth != nil {
				newCalendarAccount.OAuth2CalendarAuth = &models.OAuth2CalendarAuth{
					AccessToken:           oldCalendarAccount.GoogleCalendarAuth.AccessToken,
					RefreshToken:          oldCalendarAccount.GoogleCalendarAuth.RefreshToken,
					AccessTokenExpireDate: oldCalendarAccount.GoogleCalendarAuth.AccessTokenExpireDate,
				}
			}

			newCalendarAccounts[i] = newCalendarAccount
		}

		// Update the user document in the database
		filter := bson.M{"_id": user.Id}
		update := bson.M{"$set": bson.M{"calendarAccounts": newCalendarAccounts}}
		_, err := db.UsersCollection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			log.Printf("Error updating user %s: %v\n", user.Email, err)
		} else {
			fmt.Printf("Successfully updated user: %s\n", user.Email)
		}
	}

	fmt.Println("Migration completed.")
}
