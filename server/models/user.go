package models

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"schej.it/server/logger"
)

// Representation of a User in the mongoDB database
type User struct {
	TimezoneOffset int `json:"timezoneOffset" bson:"timezoneOffset"`

	// Profile info
	Id        primitive.ObjectID `json:"_id" bson:"_id,omitempty" binding:"required"`
	Email     string             `json:"email" bson:"email,omitempty" binding:"required"`
	FirstName string             `json:"firstName" bson:"firstName,omitempty" binding:"required"`
	LastName  string             `json:"lastName" bson:"lastName,omitempty" binding:"required"`
	Picture   string             `json:"picture" bson:"picture,omitempty" binding:"required"`

	// Calendars maps the calendar's id to the calendar object
	Calendars map[string]Calendar `json:"calendars" bson:"calendars,omitempty" binding:"required"`

	// Google OAuth stuff
	AccessToken           string             `json:"accessToken" bson:"accessToken,omitempty" binding:"required"`
	AccessTokenExpireDate primitive.DateTime `json:"accessTokenExpireDate" bson:"accessTokenExpireDate,omitempty" binding:"required"`
	RefreshToken          string             `json:"refreshToken" bson:"refreshToken,omitempty" binding:"required"`
}

// Calendar contains information about a user's calendar
type Calendar struct {
	Id       string `json:"id" bson:"id" binding:"required"`
	Summary  string `json:"summary" bson:"summary" binding:"required"`
	Selected bool   `json:"selected" bson:"selected" binding:"required"`
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
func (u *User) GetProfile() *UserProfile {
	tmp, err := json.Marshal(u)
	if err != nil {
		logger.StdErr.Panicln(err)
	}

	var profile UserProfile
	err = json.Unmarshal(tmp, &profile)
	if err != nil {
		logger.StdErr.Panicln(err)
	}
	return &profile
}
