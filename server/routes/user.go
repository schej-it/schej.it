package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"schej.it/server/db"
	"schej.it/server/middleware"
	"schej.it/server/models"
	"schej.it/server/responses"
	"schej.it/server/utils"
)

func InitUser(router *gin.Engine) {
	userRouter := router.Group("/user")
	userRouter.Use(middleware.AuthRequired())

	userRouter.GET("/profile", getProfile)
	userRouter.GET("/events", getEvents)
	userRouter.GET("/calendar", getCalendar)
}

// @Summary Gets the user's profile
// @Tags user
// @Produce json
// @Success 200 {object} models.UserProfile "A user profile object"
// @Router /user/profile [get]
func getProfile(c *gin.Context) {
	userInterface, _ := c.Get("authUser")
	user := userInterface.(*models.User)

	c.JSON(http.StatusOK, user.GetProfile())
}

// @Summary Gets all the user's events
// @Description Returns an array containing all the user's events
// @Tags user
// @Produce json
// @Success 200 {object} []models.Event
// @Router /user/events [get]
func getEvents(c *gin.Context) {
	session := sessions.Default(c)

	var events []models.Event
	cursor, err := db.EventsCollection.Find(context.Background(), bson.M{
		"ownerId": utils.GetUserId(session),
	})
	if err != nil {
		panic(err)
	}
	if err := cursor.All(context.Background(), &events); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, events)
}

// @Summary Gets the user's calendar events
// @Description Gets the user's calendar events between "timeMin" and "timeMax"
// @Tags user
// @Produce json
// @Param timeMin query string true "Lower bound for event's start time to filter by"
// @Param timeMax query string true "Upper bound for event's end time to filter by"
// @Success 200 {object} []models.CalendarEvent
// @Router /user/calendar [get]
func getCalendar(c *gin.Context) {
	// Bind query parameters
	payload := struct {
		TimeMin time.Time `form:"timeMin" binding:"required"`
		TimeMax time.Time `form:"timeMax" binding:"required"`
	}{}
	if err := c.Bind(&payload); err != nil {
		return
	}

	// Refresh token if necessary
	userInterface, _ := c.Get("authUser")
	user := userInterface.(*models.User)
	db.RefreshUserTokenIfNecessary(user)

	calendars, err := utils.GetCalendarList(user.AccessToken)
	if err != nil {
		c.JSON(err.Code, responses.Error{Error: *err})
		return
	}

	// Call the google calendar API to get a list of calendar events from the user's gcal
	// TODO: get events for all user's calendars, not just primary
	calendarEvents := make([]models.CalendarEvent, 0)
	for _, calendar := range calendars {
		events, err := utils.GetCalendarEvents(user.AccessToken, calendar.Id, payload.TimeMin, payload.TimeMax)
		if err != nil {
			c.JSON(err.Code, responses.Error{Error: *err})
			return
		}

		calendarEvents = append(calendarEvents, events...)
	}

	c.JSON(http.StatusOK, calendarEvents)
}
