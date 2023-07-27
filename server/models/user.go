package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	AccessToken           string             `json:"-" bson:"accessToken,omitempty"`
	AccessTokenExpireDate primitive.DateTime `json:"-" bson:"accessTokenExpireDate,omitempty"`
	RefreshToken          string             `json:"-" bson:"refreshToken,omitempty"`
	TokenOrigin           TokenOriginType    `json:"-" bson:"tokenOrigin,omitempty"`
}

// CalendarAccount contains info about the user's other signed in calendar accounts
type CalendarAccount struct {
	Email   string `json:"email" bson:"email,omitempty"`
	Picture string `json:"picture" bson:"picture,omitempty"`

	AccessToken           string             `json:"-" bson:"accessToken,omitempty"`
	AccessTokenExpireDate primitive.DateTime `json:"-" bson:"accessTokenExpireDate,omitempty"`
	RefreshToken          string             `json:"-" bson:"refreshToken,omitempty"`
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
