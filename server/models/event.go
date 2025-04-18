package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventType string

const (
	SPECIFIC_DATES EventType = "specific_dates"
	DOW            EventType = "dow"
	GROUP          EventType = "group"
)

// Object containing information associated with the remindee
type Remindee struct {
	Email     string   `json:"email" bson:"email,omitempty"`
	TaskIds   []string `json:"-" bson:"taskIds,omitempty"` // Task IDs of the scheduled emails
	Responded *bool    `json:"responded" bson:"responded,omitempty"`
}

type SignUpBlock struct {
	Id        primitive.ObjectID  `json:"_id" bson:"_id,omitempty"`
	Name      string              `json:"name" bson:"name,omitempty"`
	Capacity  *int                `json:"capacity" bson:"capacity,omitempty"`
	StartDate *primitive.DateTime `json:"startDate" bson:"startDate,omitempty"`
	EndDate   *primitive.DateTime `json:"endDate" bson:"endDate,omitempty"`
}

type SignUpResponse struct {
	// The IDs of the sign up blocks that the user has signed up for
	SignUpBlockIds []primitive.ObjectID `json:"signUpBlockIds" bson:"signUpBlockIds,omitempty"`

	// Guest information
	Name  string `json:"name" bson:"name,omitempty"`
	Email string `json:"email" bson:"email,omitempty"`

	// User information
	UserId primitive.ObjectID `json:"userId" bson:"userId,omitempty"`
	User   *User              `json:"user" bson:",omitempty"`
}

// Representation of an Event in the mongoDB database
type Event struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	ShortId     *string            `json:"shortId" bson:"shortId,omitempty"`
	OwnerId     primitive.ObjectID `json:"ownerId" bson:"ownerId,omitempty"`
	Name        string             `json:"name" bson:"name,omitempty"`
	Description *string            `json:"description" bson:"description,omitempty"`

	Duration                 *float32             `json:"duration" bson:"duration,omitempty"`
	Dates                    []primitive.DateTime `json:"dates" bson:"dates,omitempty"`
	NotificationsEnabled     *bool                `json:"notificationsEnabled" bson:"notificationsEnabled,omitempty"`
	SendEmailAfterXResponses *int                 `json:"sendEmailAfterXResponses" bson:"sendEmailAfterXResponses,omitempty"`
	When2meetHref            *string              `json:"when2meetHref" bson:"when2meetHref,omitempty"`
	CollectEmails            *bool                `json:"collectEmails" bson:"collectEmails,omitempty"`

	Type EventType `json:"type" bson:"type,omitempty"`

	// Sign up form details
	IsSignUpForm    *bool                      `json:"isSignUpForm" bson:"isSignUpForm,omitempty"`
	SignUpBlocks    *[]SignUpBlock             `json:"signUpBlocks" bson:"signUpBlocks,omitempty"`
	SignUpResponses map[string]*SignUpResponse `json:"signUpResponses" bson:"signUpResponses"`

	// Whether to start the event on Monday (as opposed to Sunday, used for DOW events)
	StartOnMonday *bool `json:"startOnMonday" bson:"startOnMonday,omitempty"`

	// Whether to enable blind availability
	BlindAvailabilityEnabled *bool `json:"blindAvailabilityEnabled" bson:"blindAvailabilityEnabled,omitempty"`

	// Whether to only poll for days, not times
	DaysOnly *bool `json:"daysOnly" bson:"daysOnly,omitempty"`

	// Availability responses - new format for indexed queries
	// ResponsesList []EventResponse `json:"-" bson:"responses"`
	// Availability responses - old format for backward compatibility
	ResponsesMap map[string]*Response `json:"responses" bson:"responsesMap"`

	// Scheduled event
	ScheduledEvent  *CalendarEvent `json:"scheduledEvent" bson:"scheduledEvent,omitempty"`
	CalendarEventId string         `json:"calendarEventId" bson:"calendarEventId,omitempty"`

	// Remindees
	Remindees *[]Remindee `json:"remindees" bson:"remindees,omitempty"`

	// Attendees for an availability group (fetched from Attendees collection)
	Attendees *[]Attendee `json:"attendees" bson:"-"`
}

func (e *Event) GetId() string {
	if e.ShortId != nil {
		return *e.ShortId
	}

	return e.Id.Hex()
}
