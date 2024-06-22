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

	// Whether the user has set a custom name for themselves, i.e. don't change their name when they sign in
	HasCustomName *bool `json:"hasCustomName" bson:"hasCustomName,omitempty"`

	// CalendarAccounts is a mapping from {email => CalendarAccount} that contains all the
	// additional accounts the user wants to see google calendar events for
	CalendarAccounts map[string]CalendarAccount `json:"calendarAccounts" bson:"calendarAccounts,omitempty"`

	// Google OAuth stuff
	TokenOrigin TokenOriginType `json:"-" bson:"tokenOrigin,omitempty"`

	// Calendar options
	CalendarOptions *CalendarOptions `json:"calendarOptions" bson:"calendarOptions,omitempty"`
}

// CalendarAccount contains info about the user's other signed in calendar accounts
type CalendarAccount struct {
	Email   string `json:"email" bson:"email,omitempty"`
	Picture string `json:"picture" bson:"picture,omitempty"`
	Enabled *bool  `json:"enabled" bson:"enabled,omitempty"`

	SubCalendars *map[string]SubCalendar `json:"subCalendars" bson:"subCalendars,omitempty"`

	AccessToken           string             `json:"-" bson:"accessToken,omitempty"`
	AccessTokenExpireDate primitive.DateTime `json:"-" bson:"accessTokenExpireDate,omitempty"`
	RefreshToken          string             `json:"-" bson:"refreshToken,omitempty"`
}

type SubCalendar struct {
	Name    string `json:"name" bson:"name,omitempty"`
	Enabled *bool  `json:"enabled" bson:"enabled,omitempty"`
}
type CalendarOptions struct {
	BufferTime   BufferTimeOptions   `json:"bufferTime" bson:"bufferTime"`
	WorkingHours WorkingHoursOptions `json:"workingHours" bson:"workingHours"`
}

type BufferTimeOptions struct {
	Enabled bool `json:"enabled" bson:"enabled"`
	Time    int  `json:"time" bson:"time"`
}

type WorkingHoursOptions struct {
	Enabled   bool    `json:"enabled" bson:"enabled"`
	StartTime float32 `json:"startTime" bson:"startTime"`
	EndTime   float32 `json:"endTime" bson:"endTime"`
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
