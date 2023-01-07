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
}

// @Summary Creates a new event
// @Tags events
// @Accept json
// @Produce json
// @Param payload body object{name=string,startDate=string,endDate=string} true "Object containing info about the event to create"
// @Success 201 {object} object{eventId=string}
// @Router /events [post]
func createEvent(c *gin.Context) {
	payload := struct {
		Name      string   `json:"name" binding:"required"`
		StartTime float32  `json:"startTime" binding:"required"`
		EndTime   float32  `json:"endTime" binding:"required"`
		Dates     []string `json:"dates" binding:"required"`
	}{}
	if err := c.Bind(&payload); err != nil {
		return
	}
	session := sessions.Default(c)

	event := models.Event{
		OwnerId:   utils.GetUserId(session),
		Name:      payload.Name,
		StartTime: payload.StartTime,
		EndTime:   payload.EndTime,
		Dates:     payload.Dates,
		Responses: make(map[string]models.Response),
	}

	result, err := db.EventsCollection.InsertOne(context.Background(), event)
	if err != nil {
		logger.StdErr.Panicln(err)
	}

	insertedId := result.InsertedID.(primitive.ObjectID).Hex()

	userInterface, _ := c.Get("authUser")
	user := userInterface.(*models.User)
	discord_bot.SendMessage(fmt.Sprintf(":tada: **New event created!** :tada: \n**Event name**: \"%s\"\n**Creator**: %s %s (%s)", event.Name, user.FirstName, user.LastName, user.Email))
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

	c.JSON(http.StatusOK, gin.H{})
}
