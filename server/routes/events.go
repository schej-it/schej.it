package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"schej.it/server/db"
	"schej.it/server/middleware"
	"schej.it/server/models"
	"schej.it/server/utils"
)

func InitEvents(router *gin.Engine) {
	eventRouter := router.Group("/events")

	// Creates a new event
	eventRouter.POST("", middleware.AuthRequired(), func(c *gin.Context) {
		payload := struct {
			Name      string    `json:"name" binding:"required"`
			StartDate time.Time `json:"startDate" binding:"required"`
			EndDate   time.Time `json:"endDate" binding:"required"`
			StartTime int       `json:"startTime" binding:"required"`
			EndTime   int       `json:"endTime" binding:"required"`
		}{}
		if err := c.Bind(&payload); err != nil {
			return
		}
		session := sessions.Default(c)

		event := models.Event{
			OwnerId:   utils.GetUserId(session),
			Name:      payload.Name,
			StartDate: primitive.NewDateTimeFromTime(payload.StartDate),
			EndDate:   primitive.NewDateTimeFromTime(payload.EndDate),
			StartTime: payload.StartTime,
			EndTime:   payload.EndTime,
			Responses: make(map[string]models.Response),
		}

		result, err := db.EventsCollection.InsertOne(context.Background(), event)
		if err != nil {
			panic(err)
		}

		insertedId := result.InsertedID.(primitive.ObjectID).Hex()
		c.JSON(http.StatusCreated, gin.H{"eventId": insertedId})
	})

	// Gets an event located at the given id
	eventRouter.GET("/:eventId", func(c *gin.Context) {
		eventId := c.Param("eventId")
		event := db.GetEventById(eventId)

		// Populate user fields
		for userId, response := range event.Responses {
			response.User = *db.GetUserById(userId)
			event.Responses[userId] = response
		}

		c.JSON(http.StatusOK, event)
	})

	// Updates the current user's availability
	eventRouter.POST("/:eventId/response", middleware.AuthRequired(), func(c *gin.Context) {
		payload := struct {
			Availability []string `json:"availability" binding:"required"`
		}{}
		if err := c.Bind(&payload); err != nil {
			return
		}
		session := sessions.Default(c)
		eventId := c.Param("eventId")

		response := models.Response{
			UserId:       utils.GetUserId(session),
			Availability: payload.Availability,
		}

		userIdString := session.Get("userId").(string)
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
			panic(err)
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	})
}
