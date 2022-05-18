package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"schej.it/server/db"
	"schej.it/server/middleware"
	"schej.it/server/models"
	"schej.it/server/utils"
)

func InitUser(router *gin.Engine) {
	userRouter := router.Group("/user")
	userRouter.Use(middleware.AuthRequired())

	// Gets the user's profile
	userRouter.GET("/profile", func(c *gin.Context) {
		user, _ := c.Get("authUser")
		user = user.(*models.User)

		c.JSON(http.StatusOK, user)
	})

	// Gets all the user's events
	userRouter.GET("/events", func(c *gin.Context) {
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
	})

	// Gets the times that the current user is available
	userRouter.GET("/calendar", func(c *gin.Context) {
		// Bind query parameters
		payload := struct {
			TimeMin string `form:"timeMin" binding:"required"`
			TimeMax string `form:"timeMax" binding:"required"`
		}{}
		if err := c.Bind(&payload); err != nil {
			return
		}
		session := sessions.Default(c)

		req, err := http.NewRequest(
			"GET",
			fmt.Sprintf("https://www.googleapis.com/calendar/v3/calendars/primary/events?timeMin=%s&timeMax=%s", payload.TimeMin, payload.TimeMax),
			nil,
		)
		if err != nil {
			panic(err)
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", session.Get("accessToken")))
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)

		c.JSON(http.StatusOK, res)
	})
}
