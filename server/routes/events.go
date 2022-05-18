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

	eventRouter.POST("", middleware.AuthRequired(), createEvent)
	eventRouter.GET("/:eventId", getEvent)
	eventRouter.POST("/:eventId/response", middleware.AuthRequired(), updateEventResponse)
}

// @Summary Creates a new event
// @Tags event
// @Accept json
// @Produce json
// @Param name body string true "Name"
// @Param startDate body string true "Start date"
// @Param endDate body string true "End date"
// @Param startTime body int true "Start time"
// @Param endTime body int true "End time"
// @Success 201 {object} object{eventId=string}
// @Router /events [post]
func createEvent(c *gin.Context) {
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
}

// @Summary Gets an event based on its id
// @Tags event
// @Produce json
// @Param eventId path string true "Event ID"
// @Success 200 {object} models.Event
// @Router /events/:eventId [get]
func getEvent(c *gin.Context) {
	eventId := c.Param("eventId")
	event := db.GetEventById(eventId)

	// Populate user fields
	for userId, response := range event.Responses {
		response.User = *db.GetUserById(userId)
		event.Responses[userId] = response
	}

	c.JSON(http.StatusOK, event)
}

// @Summary Updates the current user's availability
// @Tags event
// @Accept json
// @Produce json
// @Param eventId path string true "Event ID"
// @Param availability body []string true "Array of dates representing user's availability"
// @Success 200
// @Router /events/:eventId/response [post]
func updateEventResponse(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{})
}
