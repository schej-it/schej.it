/* The /user group contains all the routes to get all the information about the currently signed in user */
package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"schej.it/server/db"
	"schej.it/server/logger"
	"schej.it/server/middleware"
	"schej.it/server/models"
	"schej.it/server/responses"
	"schej.it/server/services/auth"
	"schej.it/server/services/calendar"
	"schej.it/server/services/contacts"
	"schej.it/server/utils"
)

func InitUser(router *gin.RouterGroup) {
	userRouter := router.Group("/user")
	userRouter.Use(middleware.AuthRequired())

	userRouter.GET("/profile", getProfile)
	userRouter.PATCH("/name", updateName)
	userRouter.PATCH("/calendar-options", updateCalendarOptions)
	userRouter.GET("/events", getEvents)
	userRouter.GET("/calendars", getCalendars)
	userRouter.POST("/add-google-calendar-account", addGoogleCalendarAccount)
	userRouter.DELETE("/remove-calendar-account", removeCalendarAccount)
	userRouter.POST("/toggle-calendar", toggleCalendar)
	userRouter.POST("/toggle-sub-calendar", toggleSubCalendar)
	userRouter.GET("/searchContacts", searchContacts)
	userRouter.DELETE("", deleteUser)
}

// @Summary Gets the user's profile
// @Tags user
// @Produce json
// @Success 200 {object} models.User "A user profile object"
// @Router /user/profile [get]
func getProfile(c *gin.Context) {
	userInterface, _ := c.Get("authUser")
	user := userInterface.(*models.User)

	db.UpdateDailyUserLog(user)

	c.JSON(http.StatusOK, user)
}

// @Summary Updates the user's name
// @Tags user
// @Accept json
// @Produce json
// @Param payload body object{firstName=string,lastName=string} true "Object containing the updated name"
// @Success 200
// @Router /user/name [patch]
func updateName(c *gin.Context) {
	payload := struct {
		FirstName string `json:"firstName" binding:"required"`
		LastName  string `json:"lastName" binding:"required"`
	}{}
	if err := c.BindJSON(&payload); err != nil {
		return
	}

	authUser := utils.GetAuthUser(c)

	_, err := db.UsersCollection.UpdateByID(context.Background(), authUser.Id, bson.M{
		"$set": bson.M{"firstName": payload.FirstName, "lastName": payload.LastName, "hasCustomName": true},
	})
	if err != nil {
		logger.StdErr.Panicln(err)
	}

	c.JSON(http.StatusOK, gin.H{})
}

// @Summary Updates the user's calendar options
// @Tags user
// @Accept json
// @Produce json
// @Param payload body object{bufferTime=models.BufferTimeOptions,workingHours=models.WorkingHoursOptions} true "Object containing the updated options"
// @Success 200
// @Router /user/calendar-options [patch]
func updateCalendarOptions(c *gin.Context) {
	payload := struct {
		BufferTime   *models.BufferTimeOptions   `json:"bufferTime"`
		WorkingHours *models.WorkingHoursOptions `json:"workingHours"`
	}{}
	if err := c.BindJSON(&payload); err != nil {
		return
	}

	authUser := utils.GetAuthUser(c)

	// Set default values for calendar options if nil
	if authUser.CalendarOptions == nil {
		authUser.CalendarOptions = &models.CalendarOptions{
			BufferTime: models.BufferTimeOptions{
				Enabled: false,
				Time:    15,
			},
			WorkingHours: models.WorkingHoursOptions{
				Enabled:   false,
				StartTime: 9,
				EndTime:   17,
			},
		}
	}

	// Update calendar options
	if payload.BufferTime != nil {
		authUser.CalendarOptions.BufferTime = *payload.BufferTime
	}
	if payload.WorkingHours != nil {
		authUser.CalendarOptions.WorkingHours = *payload.WorkingHours
	}

	// Update database
	_, err := db.UsersCollection.UpdateByID(context.Background(), authUser.Id, bson.M{
		"$set": bson.M{"calendarOptions": authUser.CalendarOptions},
	})
	if err != nil {
		logger.StdErr.Panicln(err)
	}

	c.JSON(http.StatusOK, gin.H{})
}

// @Summary Gets all the user's events
// @Description Returns an array containing all the user's events
// @Tags user
// @Produce json
// @Success 200 {object} object{events=[]models.Event,joinedEvents=[]models.Event}
// @Router /user/events [get]
func getEvents(c *gin.Context) {
	user := utils.GetAuthUser(c)
	userId := user.Id

	// Get the events associated with the current user
	events := make([]models.Event, 0)
	opts := options.Find().SetSort(bson.M{"_id": -1})
	cursor, err := db.EventsCollection.Find(context.Background(), bson.M{
		"$or": bson.A{
			bson.M{"ownerId": userId},
			bson.M{"responses." + userId.Hex(): bson.M{"$exists": true}},
			bson.M{"attendees": bson.M{"email": user.Email, "declined": false}},
		},
	}, opts)
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
		// Get rid of responses so we don't send too much data when fetching all events
		for userId := range event.Responses {
			event.Responses[userId] = nil
		}

		// Filter into events user created and responded to
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
// @Param accounts query string false "Comma separated list of accounts to fetch calendar events from"
// @Success 200 {object} map[string]calendar.CalendarEventsWithError
// @Router /user/calendars [get]
func getCalendars(c *gin.Context) {
	// Bind query parameters
	payload := struct {
		TimeMin  time.Time `form:"timeMin" binding:"required"`
		TimeMax  time.Time `form:"timeMax" binding:"required"`
		Accounts string    `form:"accounts"`
	}{}
	if err := c.Bind(&payload); err != nil {
		return
	}

	var accounts []string
	if len(payload.Accounts) == 0 {
		accounts = make([]string, 0)
	} else {
		accounts = utils.ParseArrayQueryParam(payload.Accounts)
	}
	accountsSet := utils.ArrayToSet(accounts)
	user := utils.GetAuthUser(c)

	calendarEvents, editedCalendarAccounts := calendar.GetUsersCalendarEvents(user, accountsSet, payload.TimeMin, payload.TimeMax)

	if editedCalendarAccounts {
		db.UsersCollection.FindOneAndUpdate(
			context.Background(),
			bson.M{"_id": user.Id},
			bson.M{"$set": user},
		)
	}

	c.JSON(http.StatusOK, calendarEvents)
}

// @Summary Adds a new calendar account
// @Tags user
// @Accept json
// @Produce json
// @Param payload body object{code=string} true "Object containing the Google authorization code"
// @Success 200
// @Router /user/add-google-calendar-account [post]
func addGoogleCalendarAccount(c *gin.Context) {
	payload := struct {
		Code string `json:"code" binding:"required"`
	}{}
	if err := c.BindJSON(&payload); err != nil {
		return
	}

	// Get auth user
	authUser := utils.GetAuthUser(c)

	// Get tokens
	tokens := auth.GetTokensFromAuthCode(payload.Code, utils.GetOrigin(c))

	// Get user info from JWT
	claims := utils.ParseJWT(tokens.IdToken)
	email, _ := claims.GetStr("email")
	picture, _ := claims.GetStr("picture")

	// Get access token expire time
	accessTokenExpireDate := utils.GetAccessTokenExpireDate(tokens.ExpiresIn)

	// Define a new calendar account
	calendarAccount := models.CalendarAccount{
		CalendarType: models.GoogleCalendarType,
		Details: models.GoogleCalendar{
			Email:   email,
			Picture: picture,

			AccessToken:           tokens.AccessToken,
			AccessTokenExpireDate: primitive.NewDateTimeFromTime(accessTokenExpireDate),
			RefreshToken:          tokens.RefreshToken,
		},
		Enabled: &[]bool{true}[0], // Workaround to pass a boolean pointer
	}
	calendarAccountKey := utils.GetCalendarAccountKey(email, models.GoogleCalendarType)

	// Set subcalendars map based on whether calendar account already exists
	if oldCalendarAccount, ok := authUser.CalendarAccounts[calendarAccountKey]; ok && oldCalendarAccount.SubCalendars != nil {
		calendarAccount.SubCalendars = oldCalendarAccount.SubCalendars
	} else {
		subCalendars, err := calendarAccount.Details.(calendar.CalendarProvider).GetCalendarList()
		if err == nil {
			calendarAccount.SubCalendars = &subCalendars
		}
	}

	// Set calendar account
	authUser.CalendarAccounts[calendarAccountKey] = calendarAccount

	// Perform mongo update
	db.UsersCollection.FindOneAndUpdate(
		context.Background(),
		bson.M{"_id": authUser.Id},
		bson.M{"$set": authUser},
	)

	c.JSON(http.StatusOK, gin.H{})
}

// @Summary Removes an existing calendar account
// @Tags user
// @Accept json
// @Produce json
// @Param payload body object{email=string,calendarType=models.CalendarType} true "Object containing the email + type of the calendar account to remove"
// @Success 200
// @Router /user/remove-calendar-account [delete]
func removeCalendarAccount(c *gin.Context) {
	payload := struct {
		Email        string              `json:"email" binding:"required"`
		CalendarType models.CalendarType `json:"calendarType" binding:"required"`
	}{}
	if err := c.BindJSON(&payload); err != nil {
		return
	}

	calendarAccountKey := utils.GetCalendarAccountKey(payload.Email, payload.CalendarType)

	authUser := utils.GetAuthUser(c)
	db.UsersCollection.UpdateByID(context.Background(), authUser.Id, bson.A{
		bson.M{"$set": bson.M{
			"calendarAccounts": bson.M{
				"$setField": bson.M{
					"field": calendarAccountKey,
					"input": "$$ROOT.calendarAccounts",
					"value": "$$REMOVE",
				},
			},
		}},
	})

	c.JSON(http.StatusOK, gin.H{})
}

// @Summary Toggles whether the specified calendar is enabled or disabled for the user
// @Tags user
// @Accept json
// @Produce json
// @Param payload body object{email=string,enabled=bool} true "Email of calendar account and whether to enable it"
// @Success 200
// @Router /user/toggle-calendar [post]
func toggleCalendar(c *gin.Context) {
	payload := struct {
		Email   string `json:"email" binding:"required"`
		Enabled *bool  `json:"enabled" binding:"required"`
	}{}
	if err := c.Bind(&payload); err != nil {
		logger.StdErr.Panicln(err)
		return
	}

	// Update enabled status for the specified account
	authUser := utils.GetAuthUser(c)
	if account, ok := authUser.CalendarAccounts[payload.Email]; ok {
		account.Enabled = payload.Enabled
		authUser.CalendarAccounts[payload.Email] = account

		_, err := db.UsersCollection.UpdateOne(context.Background(), bson.M{
			"_id": authUser.Id,
		}, bson.M{
			"$set": authUser,
		})
		if err != nil {
			logger.StdErr.Panicln(err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{})
}

// @Summary Toggles whether the specified sub-calendar is enabled or disabled for the user
// @Tags user
// @Accept json
// @Produce json
// @Param payload body object{email=string,subCalendarId=string,enabled=bool} true "Email of calendar account, the sub calendar id, and whether to enable it"
// @Success 200
// @Router /user/toggle-sub-calendar [post]
func toggleSubCalendar(c *gin.Context) {
	payload := struct {
		Email         string `json:"email" binding:"required"`
		SubCalendarId string `json:"subCalendarId" binding:"required"`
		Enabled       *bool  `json:"enabled" binding:"required"`
	}{}
	if err := c.Bind(&payload); err != nil {
		logger.StdErr.Panicln(err)
		return
	}

	// Update enabled status for the specified sub calendar
	authUser := utils.GetAuthUser(c)
	if account, ok := authUser.CalendarAccounts[payload.Email]; ok {
		if subCalendar, ok := (*account.SubCalendars)[payload.SubCalendarId]; ok {
			subCalendar.Enabled = payload.Enabled
			(*account.SubCalendars)[payload.SubCalendarId] = subCalendar
			authUser.CalendarAccounts[payload.Email] = account

			_, err := db.UsersCollection.UpdateOne(context.Background(), bson.M{
				"_id": authUser.Id,
			}, bson.M{
				"$set": authUser,
			})
			if err != nil {
				logger.StdErr.Panicln(err)
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{})
}

// @Summary Searches the user's contacts based on the given query
// @Tags user
// @Produce json
// @Param query query string true "Query to search for"
// @Success 200 {object} []models.User
// @Router /user/searchContacts [get]
func searchContacts(c *gin.Context) {
	// Bind query parameters
	payload := struct {
		Query string `form:"query"`
	}{}
	if err := c.Bind(&payload); err != nil {
		return
	}

	userInterface, _ := c.Get("authUser")
	user := userInterface.(*models.User)

	contacts, googleError := contacts.SearchContacts(user, payload.Query)
	if googleError != nil {
		c.JSON(googleError.Code, responses.Error{Error: *googleError})
		return
	}

	c.JSON(http.StatusOK, contacts)
}

// @Summary Deletes the currently signed in user
// @Tags user
// @Produce json
// @Success 200
// @Router /user [delete]
func deleteUser(c *gin.Context) {
	userInterface, _ := c.Get("authUser")
	user := userInterface.(*models.User)

	_, err := db.UsersCollection.DeleteOne(context.Background(), bson.M{"_id": user.Id})
	if err != nil {
		logger.StdErr.Panicln(err)
	}

	// Delete session
	session := sessions.Default(c)
	session.Delete("userId")
	session.Save()

	c.JSON(http.StatusOK, gin.H{})
}
