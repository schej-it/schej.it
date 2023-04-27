/* The /events group contains all the routes to get and edit events */
package routes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"schej.it/server/db"
	"schej.it/server/discord_bot"
	"schej.it/server/errs"
	"schej.it/server/logger"
	"schej.it/server/middleware"
	"schej.it/server/models"
	"schej.it/server/responses"
	"schej.it/server/utils"
)

func InitEvents(router *gin.Engine) {
	eventRouter := router.Group("/events")

	eventRouter.POST("", middleware.AuthRequired(), createEvent)
	eventRouter.GET("/:eventId", getEvent)
	eventRouter.POST("/:eventId/response", updateEventResponse)
	eventRouter.POST("/:eventId/schedule", middleware.AuthRequired(), scheduleEvent)
	eventRouter.PUT("/:eventId", middleware.AuthRequired(), editEvent)
	eventRouter.DELETE("/:eventId", middleware.AuthRequired(), deleteEvent)
}

// @Summary Creates a new event
// @Tags events
// @Accept json
// @Produce json
// @Param payload body object{name=string,duration=*float32,dates=[]primitive.DateTime} true "Object containing info about the event to create"
// @Success 201 {object} object{eventId=string}
// @Router /events [post]
func createEvent(c *gin.Context) {
	payload := struct {
		Name                 string               `json:"name" binding:"required"`
		Duration             *float32             `json:"duration" binding:"required"`
		Dates                []primitive.DateTime `json:"dates" binding:"required"`
		NotificationsEnabled *bool                `json:"notificationsEnabled" binding:"required"`
	}{}
	if err := c.Bind(&payload); err != nil {
		return
	}
	session := sessions.Default(c)

	event := models.Event{
		OwnerId:              utils.GetUserId(session),
		Name:                 payload.Name,
		Duration:             payload.Duration,
		Dates:                payload.Dates,
		NotificationsEnabled: *payload.NotificationsEnabled,
		Responses:            make(map[string]*models.Response),
	}

	result, err := db.EventsCollection.InsertOne(context.Background(), event)
	if err != nil {
		logger.StdErr.Panicln(err)
	}

	insertedId := result.InsertedID.(primitive.ObjectID).Hex()

	userInterface, _ := c.Get("authUser")
	user := userInterface.(*models.User)
	discord_bot.SendMessage(fmt.Sprintf(":tada: **New event created!** :tada: \n**Event url**: https://schej.it/e/%s\n**Creator**: %s %s (%s)\n**Notifications Enabled**: %v", insertedId, user.FirstName, user.LastName, user.Email, event.NotificationsEnabled))
	c.JSON(http.StatusCreated, gin.H{"eventId": insertedId})
}

// @Summary Gets an event based on its id
// @Tags events
// @Produce json
// @Param eventId path string true "Event ID"
// @Success 200 {object} models.Event
// @Router /events/{eventId} [get]
func getEvent(c *gin.Context) {
	eventId := c.Param("eventId")
	event := db.GetEventById(eventId)

	if event == nil {
		c.JSON(http.StatusNotFound, responses.Error{Error: errs.EventNotFound})
		return
	}

	// Populate user fields
	for userId, response := range event.Responses {
		user := db.GetUserById(userId)
		if user == nil {
			userId = response.Name
			response.User = &models.UserProfile{
				FirstName: response.Name,
			}
		} else {
			response.User = user.GetProfile()
		}
		event.Responses[userId] = response
	}

	c.JSON(http.StatusOK, event)
}

// @Summary Updates the current user's availability
// @Tags events
// @Accept json
// @Produce json
// @Param eventId path string true "Event ID"
// @Param payload body object{availability=[]string,guest=bool,name=string} true "Object containing info about the event response to update"
// @Success 200
// @Router /events/{eventId}/response [post]
func updateEventResponse(c *gin.Context) {
	payload := struct {
		Availability []string `json:"availability" binding:"required"`
		Guest        *bool    `json:"guest" binding:"required"`
		Name         string   `json:"name"`
	}{}
	if err := c.Bind(&payload); err != nil {
		return
	}
	session := sessions.Default(c)
	eventId := c.Param("eventId")
	event := db.GetEventById(eventId)

	var response models.Response
	var userIdString string
	// Populate response differently if guest vs signed in user
	if *payload.Guest {
		userIdString = payload.Name

		response = models.Response{
			Name:         payload.Name,
			Availability: payload.Availability,
		}
	} else {
		userIdInterface := session.Get("userId")
		if userIdInterface == nil {
			c.JSON(http.StatusUnauthorized, responses.Error{Error: errs.NotSignedIn})
			c.Abort()
			return
		}
		userIdString = userIdInterface.(string)

		response = models.Response{
			UserId:       utils.StringToObjectID(userIdString),
			Availability: payload.Availability,
		}
	}

	// Check if user has responded to event before (edit response) or not (new response)
	_, userHasResponded := event.Responses[userIdString]

	// Update responses in mongodb
	_, err := db.EventsCollection.UpdateByID(
		context.Background(),
		utils.StringToObjectID(eventId),
		bson.M{
			"$set": bson.M{
				"responses." + userIdString: response,
			},
		},
	)
	if err != nil {
		logger.StdErr.Panicln(err)
	}

	// Send email to creator of event if creator enabled it
	if event.NotificationsEnabled && !userHasResponded && userIdString != event.OwnerId.Hex() {
		// Send email asynchronously
		go func() {
			creator := db.GetUserById(event.OwnerId.Hex())
			if creator == nil {
				c.JSON(http.StatusOK, gin.H{})
				return
			}

			var respondentName string
			if *payload.Guest {
				respondentName = payload.Name
			} else {
				respondent := db.GetUserById(userIdString)
				respondentName = fmt.Sprintf("%s %s", respondent.FirstName, respondent.LastName)
			}
			utils.SendEmail(
				creator.Email,
				fmt.Sprintf("Someone just responded to your schej - \"%s\"!", event.Name),
				fmt.Sprintf(
					`<p>Hi %s,</p>

					<p>%s just responded to your schej named "%s"!<br>
					<a href="https://schej.it/e/%s">Click here to view the event</a></p>

					<p>Best,<br>
					schej team</p>`,
					creator.FirstName, respondentName, event.Name, eventId,
				),
				"text/html",
			)
		}()
	}

	c.JSON(http.StatusOK, gin.H{})
}

// @Summary Schedules an event on the user's google calendar
// @Tags events
// @Accept json
// @Produce json
// @Param eventId path string true "Event ID"
// @Param payload body object{startDate=primitive.DateTime,endDate=primitive.DateTime,attendeeEmails=[]string} true "Object containing info about the event to schedule"
// @Success 200
// @Router /events/{eventId}/schedule [post]
func scheduleEvent(c *gin.Context) {
	payload := struct {
		StartDate *primitive.DateTime `json:"startDate" binding:"required"`
		EndDate   *primitive.DateTime `json:"endDate" binding:"required"`
		AttendeeEmails []string `json:"attendeeEmails" binding:"required"`
	}{}
	if err := c.Bind(&payload); err != nil {
		return
	}
	eventId := c.Param("eventId")
	event := db.GetEventById(eventId)

	// TODO: update event if calendarEventId exists

	// Create google calendar invite
	userInterface, _ := c.Get("authUser")
	user := userInterface.(*models.User)
	calendarEventId, googleApiError := db.ScheduleEvent(user, event.Name, *payload.StartDate, *payload.EndDate, payload.AttendeeEmails)
	if googleApiError != nil {
		c.JSON(googleApiError.Code, responses.Error{Error: *googleApiError})
		return
	}

	scheduledEvent := models.CalendarEvent{
		StartDate: *payload.StartDate,
		EndDate: *payload.EndDate,
	}

	// Update event
	_, err := db.EventsCollection.UpdateByID(
		context.Background(),
		utils.StringToObjectID(eventId),
		bson.M{
			"$set": bson.M{
				"scheduledEvent": scheduledEvent,
				"calendarEventId": calendarEventId,
			},
		},
	)
	if err != nil {
		logger.StdErr.Panicln(err)
	}


	c.JSON(http.StatusOK, gin.H{})
}

// @Summary Edits an event based on its id
// @Tags events
// @Produce json
// @Param eventId path string true "Event ID"
// @Param payload body object{name=string,duration=*float32,dates=[]primitive.DateTime} true "Object containing info about the event to update"
// @Success 200
// @Router /events/{eventId} [put]
func editEvent(c *gin.Context) {
	payload := struct {
		Name                 string               `json:"name" binding:"required"`
		Duration             *float32             `json:"duration" binding:"required"`
		Dates                []primitive.DateTime `json:"dates" binding:"required"`
		NotificationsEnabled *bool                `json:"notificationsEnabled" binding:"required"`
	}{}
	if err := c.Bind(&payload); err != nil {
		return
	}

	eventId := c.Param("eventId")
	objectId, err := primitive.ObjectIDFromHex(eventId)
	if err != nil {
		// eventId is malformatted
		c.Status(http.StatusBadRequest)
		return
	}

	userInterface, _ := c.Get("authUser")
	user := userInterface.(*models.User)

	_, err = db.EventsCollection.UpdateOne(
		context.Background(),
		bson.M{
			"_id":     objectId,
			"ownerId": user.Id,
		},
		bson.M{
			"$set": bson.M{
				"name":                 payload.Name,
				"duration":             payload.Duration,
				"dates":                payload.Dates,
				"notificationsEnabled": *payload.NotificationsEnabled,
			},
		},
	)

	if err != nil {
		logger.StdErr.Panicln(err)
	}

	c.Status(http.StatusOK)
}

// @Summary Deletes an event based on its id
// @Tags events
// @Produce json
// @Param eventId path string true "Event ID"
// @Success 200
// @Router /events/{eventId} [delete]
func deleteEvent(c *gin.Context) {
	eventId := c.Param("eventId")

	objectId, err := primitive.ObjectIDFromHex(eventId)
	if err != nil {
		// eventId is malformatted
		c.Status(http.StatusBadRequest)
		return
	}

	userInterface, _ := c.Get("authUser")
	user := userInterface.(*models.User)

	_, err = db.EventsCollection.DeleteOne(context.Background(), bson.M{
		"_id":     objectId,
		"ownerId": user.Id,
	})
	if err != nil {
		logger.StdErr.Panicln(err)
	}

	c.Status(http.StatusOK)
}
