package models

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Representation of a User in the mongoDB database
type User struct {
	// Profile info
	Id        primitive.ObjectID `json:"_id" bson:"_id,omitempty" binding:"required"`
	Email     string             `json:"email" bson:"email" binding:"required"`
	FirstName string             `json:"firstName" bson:"firstName" binding:"required"`
	LastName  string             `json:"lastName" bson:"lastName" binding:"required"`
	Picture   string             `json:"picture" bson:"picture" binding:"required"`

	// Google OAuth stuff
	AccessToken           string             `json:"accessToken" binding:"required"`
	AccessTokenExpireDate primitive.DateTime `json:"accessTokenExpireDate" bson:"accessTokenExpireDate" binding:"required"`
	RefreshToken          string             `json:"refreshToken" bson:"refreshToken" binding:"required"`
}

// User profile to return as json to frontend
type UserProfile struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id,omitempty" binding:"required"`
	Email     string             `json:"email" bson:"email" binding:"required"`
	FirstName string             `json:"firstName" bson:"firstName" binding:"required"`
	LastName  string             `json:"lastName" bson:"lastName" binding:"required"`
	Picture   string             `json:"picture" bson:"picture" binding:"required"`
}

// Get a UserProfile object from the given User object
func (u *User) GetProfile() UserProfile {

	tmp, err := json.Marshal(u)
	if err != nil {
		panic(err)
	}

	var profile UserProfile
	err = json.Unmarshal(tmp, &profile)
	if err != nil {
		panic(err)
	}
	return profile
}

// If access token has expired, get a new token, update the user object, and save it to the database
/*func (u *User) RefreshTokenIfNecessary() {
	if time.Now().After(u.AccessTokenExpireDate.Time()) {
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

		db.UsersCollection.FindOneAndUpdate(
			context.Background(),
			bson.M{"email": u.Email},
			bson.M{"$set": u},
		)
	}
}
*/
