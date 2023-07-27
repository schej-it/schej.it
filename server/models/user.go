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

	// CalendarAccounts contains all the additional accounts the user wants to see google calendar events for
	CalendarAccounts []CalendarAccount `json:"calendarAccounts" bson:"calendarAccounts,omitempty"`

	// Google OAuth stuff
	AccessToken           string             `json:"accessToken" bson:"accessToken,omitempty"`
	AccessTokenExpireDate primitive.DateTime `json:"accessTokenExpireDate" bson:"accessTokenExpireDate,omitempty"`
	RefreshToken          string             `json:"refreshToken" bson:"refreshToken,omitempty"`
	TokenOrigin           TokenOriginType    `json:"tokenOrigin" bson:"tokenOrigin,omitempty"`
}

// CalendarAccount contains info about the user's other signed in calendar accounts
type CalendarAccount struct {
	Email   string `json:"email" bson:"email,omitempty"`
	Picture string `json:"picture" bson:"picture,omitempty"`

	AccessToken           string             `json:"-" bson:"accessToken,omitempty"`
	AccessTokenExpireDate primitive.DateTime `json:"-" bson:"accessTokenExpireDate,omitempty"`
	RefreshToken          string             `json:"-" bson:"refreshToken,omitempty"`
}

// User profile to return as json to frontend
type UserProfile struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Email     string             `json:"email" bson:"email,omitempty"`
	FirstName string             `json:"firstName" bson:"firstName,omitempty"`
	LastName  string             `json:"lastName" bson:"lastName,omitempty"`
	Picture   string             `json:"picture" bson:"picture,omitempty"`

	CalendarAccounts []CalendarAccount `json:"calendarAccounts" bson:"calendarAccounts,omitempty"`
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

type UserStatus string

const (
	FREE UserStatus = "free"
	BUSY UserStatus = "busy"
)
