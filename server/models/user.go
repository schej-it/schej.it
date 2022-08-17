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
	Id        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Email     string             `json:"email" bson:"email,omitempty"`
	FirstName string             `json:"firstName" bson:"firstName,omitempty"`
	LastName  string             `json:"lastName" bson:"lastName,omitempty"`
	Picture   string             `json:"picture" bson:"picture,omitempty"`

	// Settings
	Visibility int `json:"visibility" bson:"visibility"`

	// Friends
	FriendIds []primitive.ObjectID `json:"-" bson:"friendIds"`
	Friends   []UserProfile        `json:"friends" bson:",omitempty"`

	// Calendars maps the calendar's id to the calendar object
	Calendars map[string]Calendar `json:"calendars" bson:"calendars"`

	// Google OAuth stuff
	AccessToken           string             `json:"accessToken" bson:"accessToken,omitempty"`
	AccessTokenExpireDate primitive.DateTime `json:"accessTokenExpireDate" bson:"accessTokenExpireDate,omitempty"`
	RefreshToken          string             `json:"refreshToken" bson:"refreshToken,omitempty"`
	TokenOrigin           TokenOriginType    `json:"tokenOrigin" bson:"tokenOrigin,omitempty"`
}

// Calendar contains information about a user's calendar
type Calendar struct {
	Id       string `json:"id" bson:"id,omitempty"`
	Summary  string `json:"summary" bson:"summary,omitempty"`
	Selected bool   `json:"selected" bson:"selected,omitempty"`
}

// User profile to return as json to frontend
type UserProfile struct {
	Id         primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Email      string             `json:"email" bson:"email,omitempty"`
	FirstName  string             `json:"firstName" bson:"firstName,omitempty"`
	LastName   string             `json:"lastName" bson:"lastName,omitempty"`
	Picture    string             `json:"picture" bson:"picture,omitempty"`
	Visibility *int               `json:"visibility" bson:"visibility,omitempty"`
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

// Declare the possible types of TokenOrigin
type TokenOriginType string

const (
	Undefined TokenOriginType = ""
	IOS       TokenOriginType = "ios"
	ANDROID   TokenOriginType = "android"
	WEB       TokenOriginType = "web"
)
