package db

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"schej.it/server/models"
	"schej.it/server/utils"
)

func GetUserById(userId string) *models.User {
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		// eventId is malformatted
		return nil
	}
	result := UsersCollection.FindOne(context.Background(), bson.M{
		"_id": objectId,
	})
	if result.Err() == mongo.ErrNoDocuments {
		// User does not exist!
		return nil
	}

	// Set auth user request variable
	var user models.User
	if err := result.Decode(&user); err != nil {
		panic(err)
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

	// Set auth user request variable
	var event models.Event
	if err := result.Decode(&event); err != nil {
		panic(err)
	}

	return &event
}

// If access token has expired, get a new token, update the user object, and save it to the database
func RefreshUserTokenIfNecessary(u *models.User) {
	if time.Now().After(u.AccessTokenExpireDate.Time()) {
		fmt.Println("REFRESHING TOKEN")

		// Refresh token by calling google token endpoint
		values := url.Values{
			"client_id":     {os.Getenv("CLIENT_ID")},
			"client_secret": {os.Getenv("CLIENT_SECRET")},
			"grant_type":    {"refresh_token"},
			"refresh_token": {u.RefreshToken},
		}
		resp, err := http.PostForm(
			"https://oauth2.googleapis.com/token",
			values,
		)
		if err != nil {
			panic(err)
		}
		res := struct {
			AccessToken string `json:"access_token"`
			ExpiresIn   int    `json:"expires_in"`
			Scope       string `json:"scope"`
			TokenType   string `json:"token_type"`
		}{}
		json.NewDecoder(resp.Body).Decode(&res)

		accessTokenExpireDate := utils.GetAccessTokenExpireDate(res.ExpiresIn)
		u.AccessToken = res.AccessToken
		u.AccessTokenExpireDate = primitive.NewDateTimeFromTime(accessTokenExpireDate)

		UsersCollection.FindOneAndUpdate(
			context.Background(),
			bson.M{"email": u.Email},
			bson.M{"$set": u},
		)
	}
}
