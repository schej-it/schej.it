/* The /user group contains all the routes to get all the information about the currently signed in user */
package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"schej.it/server/db"
	"schej.it/server/logger"
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
	userRouter.POST("/visibility", updateVisibility)
}

// @Summary Gets the user's profile
// @Tags user
// @Produce json
// @Success 200 {object} models.UserProfile "A user profile object"
// @Router /user/profile [get]
func getProfile(c *gin.Context) {
	userInterface, _ := c.Get("authUser")
	user := userInterface.(*models.User)

	db.UpdateDailyUserLog(user)

	c.JSON(http.StatusOK, user.GetProfile())
}

// @Summary Gets all the user's events
// @Description Returns an array containing all the user's events
// @Tags user
// @Produce json
// @Success 200 {object} object{events=[]models.Event,joinedEvents=[]models.Event}
// @Router /user/events [get]
func getEvents(c *gin.Context) {
	session := sessions.Default(c)

	userId := utils.GetUserId(session)
	userIdString := session.Get("userId").(string)

	// Get the events associated with the current user
	events := make([]models.Event, 0)
	cursor, err := db.EventsCollection.Find(context.Background(), bson.M{
		"$or": bson.A{
			bson.M{"ownerId": userId},
			bson.M{"responses." + userIdString: bson.M{"$exists": true}},
		},
	})
	if err != nil {
		logger.StdErr.Panicln(err)
	}
	if err := cursor.All(context.Background(), &events); err != nil {
		logger.StdErr.Panicln(err)
	}

	response := make(map[string][]models.Event)
	response["events"] = make([]models.Event, 0)       // The events the user created
	response["joinedEvents"] = make([]models.Event, 0) // The events the user has responded to

	for _, event := range events {
		if event.OwnerId == userId {
			response["events"] = append(response["events"], event)
		} else {
			response["joinedEvents"] = append(response["joinedEvents"], event)
		}
	}

	c.JSON(http.StatusOK, response)
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

// @Summary Updates the current user's visibility
// @Tags user
// @Accept json
// @Produce json
// @Param payload body object{visibility=int} true "Visibility of user from 0 to 2"
// @Success 200
// @Router /user/visibility [post]
func updateVisibility(c *gin.Context) {

	// Bind query parameters
	payload := struct {
		Visibility *int `json:"visibility" binding:"required"`
	}{}
	if err := c.Bind(&payload); err != nil {
		logger.StdErr.Panicln(err)
		return
	}

	session := sessions.Default(c)
	userId := utils.GetUserId(session)

	_, err := db.UsersCollection.UpdateByID(
		context.Background(),
		userId,
		bson.M{
			"$set": bson.M{
				"visibility": payload.Visibility,
			},
		},
	)
	if err != nil {
		logger.StdErr.Panicln(err)
	}

	c.JSON(http.StatusOK, gin.H{})
}
