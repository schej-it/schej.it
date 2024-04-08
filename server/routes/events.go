/* The /events group contains all the routes to get and edit events */
package routes

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"schej.it/server/db"
	"schej.it/server/errs"
	"schej.it/server/logger"
	"schej.it/server/middleware"
	"schej.it/server/models"
	"schej.it/server/responses"
	"schej.it/server/services/calendar"
	"schej.it/server/services/gcloud"
	"schej.it/server/services/listmonk"
	"schej.it/server/slackbot"
	"schej.it/server/utils"
)

func InitEvents(router *gin.Engine) {
	eventRouter := router.Group("/events")

	eventRouter.POST("", createEvent)
	eventRouter.PUT("/:eventId", editEvent)
	eventRouter.GET("/:eventId", getEvent)
	eventRouter.POST("/:eventId/response", updateEventResponse)
	eventRouter.DELETE("/:eventId/response", deleteEventResponse)
	eventRouter.POST("/:eventId/responded", userResponded)
	eventRouter.POST("/:eventId/decline", middleware.AuthRequired(), declineInvite)
	eventRouter.GET("/:eventId/calendar-availabilities", middleware.AuthRequired(), getCalendarAvailabilities)
	eventRouter.DELETE("/:eventId", middleware.AuthRequired(), deleteEvent)
	eventRouter.POST("/:eventId/duplicate", middleware.AuthRequired(), duplicateEvent)
}

// @Summary Creates a new event
// @Tags events
// @Accept json
// @Produce json
// @Param payload body object{name=string,duration=float32,dates=[]string,type=models.EventType,notificationsEnabled=bool,remindees=[]string,attendees=[]string} true "Object containing info about the event to create"
// @Success 201 {object} object{eventId=string}
// @Router /events [post]
func createEvent(c *gin.Context) {
	payload := struct {
		// Required parameters
		Name     string               `json:"name" binding:"required"`
		Duration *float32             `json:"duration" binding:"required"`
		Dates    []primitive.DateTime `json:"dates" binding:"required"`
		Type     models.EventType     `json:"type" binding:"required"`

		// Only for discrete events
		NotificationsEnabled *bool    `json:"notificationsEnabled"`
		Remindees            []string `json:"remindees"`

		// Only for availability groups
		Attendees []string `json:"attendees"`
	}{}
	if err := c.Bind(&payload); err != nil {
		return
	}
	session := sessions.Default(c)

	// If user logged in, set owner id to their user id, otherwise set owner id to nil
	userIdInterface := session.Get("userId")
	userId, signedIn := userIdInterface.(string)
	var user *models.User
	var ownerId primitive.ObjectID
	if signedIn {
		ownerId = utils.StringToObjectID(userId)
		user = db.GetUserById(userId)
	} else {
		ownerId = primitive.NilObjectID
	}

	// Construct event object
	event := models.Event{
		Id:                   primitive.NewObjectID(),
		OwnerId:              ownerId,
		Name:                 payload.Name,
		Duration:             payload.Duration,
		Dates:                payload.Dates,
		NotificationsEnabled: payload.NotificationsEnabled,
		Type:                 payload.Type,
		Responses:            make(map[string]*models.Response),
	}

	// Generate short id
	shortId := db.GenerateShortEventId(event.Id)
	event.ShortId = &shortId

	// Schedule reminder emails if remindees array is not empty
	if len(payload.Remindees) > 0 {
		// Determine owner name
		var ownerName string
		if signedIn {
			ownerName = user.FirstName
		} else {
			ownerName = "Somebody"
		}

		// Schedule email reminders for each of the remindees' emails
		remindees := make([]models.Remindee, 0)
		for _, email := range payload.Remindees {
			taskIds := gcloud.CreateEmailTask(email, ownerName, payload.Name, event.GetId())
			remindees = append(remindees, models.Remindee{
				Email:     email,
				TaskIds:   taskIds,
				Responded: utils.FalsePtr(),
			})
		}

		event.Remindees = &remindees
	}

	if payload.Type == models.GROUP {
		attendees := make([]models.Attendee, 0)

		if signedIn {
			// Add event owner to group by default
			enabledCalendars := make(map[string][]string)
			for email, calendarAccount := range user.CalendarAccounts {
				if utils.Coalesce(calendarAccount.Enabled) {
					enabledCalendars[email] = make([]string, 0)
					for calendarId, subCalendar := range utils.Coalesce(calendarAccount.SubCalendars) {
						if utils.Coalesce(subCalendar.Enabled) {
							enabledCalendars[email] = append(enabledCalendars[email], calendarId)
						}
					}
				}
			}
			event.Responses[user.Id.Hex()] = &models.Response{
				UserId:                  user.Id,
				UseCalendarAvailability: utils.TruePtr(),
				EnabledCalendars:        &enabledCalendars,
			}

			// Add owner as attendee
			attendees = append(attendees, models.Attendee{Email: user.Email, Declined: utils.FalsePtr()})
		}

		// Add attendees and send email
		if len(payload.Attendees) > 0 {
			// Determine owner name
			var ownerName string
			if signedIn {
				ownerName = user.FirstName
			} else {
				ownerName = "Somebody"
			}

			// Add attendees to attendees array and send invite emails
			availabilityGroupInviteEmailId := 9
			for _, email := range payload.Attendees {
				listmonk.SendEmailAddSubscriberIfNotExist(email, availabilityGroupInviteEmailId, bson.M{
					"ownerName": ownerName,
					"groupName": event.Name,
					"groupUrl":  fmt.Sprintf("%s/g/%s", utils.GetBaseUrl(), event.GetId()),
				})
				attendees = append(attendees, models.Attendee{Email: email, Declined: utils.FalsePtr()})
			}

		}

		event.Attendees = &attendees
	}

	// Insert event
	result, err := db.EventsCollection.InsertOne(context.Background(), event)
	if err != nil {
		logger.StdErr.Panicln(err)
	}
	insertedId := result.InsertedID.(primitive.ObjectID).Hex()

	// Send slackbot message
	var creator string
	if signedIn {
		creator = fmt.Sprintf("%s %s (%s)", user.FirstName, user.LastName, user.Email)
	} else {
		creator = "Guest :face_with_open_eyes_and_hand_over_mouth:"
	}
	slackbot.SendEventCreatedMessage(insertedId, creator, event)

	c.JSON(http.StatusCreated, gin.H{"eventId": insertedId, "shortId": event.ShortId})
}

// @Summary Edits an event based on its id
// @Tags events
// @Produce json
// @Param eventId path string true "Event ID"
// @Param payload body object{name=string,duration=float32,dates=[]string,type=models.EventType,notificationsEnabled=bool,remindees=[]string,attendees=[]string} true "Object containing info about the event to update"
// @Success 200
// @Router /events/{eventId} [put]
func editEvent(c *gin.Context) {
	payload := struct {
		// Required parameters
		Name     string               `json:"name" binding:"required"`
		Duration *float32             `json:"duration" binding:"required"`
		Dates    []primitive.DateTime `json:"dates" binding:"required"`
		Type     models.EventType     `json:"type" binding:"required"`

		// Only for discrete events
		NotificationsEnabled *bool    `json:"notificationsEnabled"`
		Remindees            []string `json:"remindees"`

		// Only for availability groups
		Attendees []string `json:"attendees"`
	}{}
	if err := c.Bind(&payload); err != nil {
		return
	}

	eventId := c.Param("eventId")
	event := db.GetEventByEitherId(eventId)
	if event == nil {
		c.JSON(http.StatusNotFound, responses.Error{Error: errs.EventNotFound})
		return
	}

	// If user logged in, set owner id to their user id, otherwise set owner id to nil
	session := sessions.Default(c)
	userIdInterface := session.Get("userId")
	userId, signedIn := userIdInterface.(string)
	var ownerId primitive.ObjectID
	if signedIn {
		ownerId = utils.StringToObjectID(userId)
	} else {
		ownerId = primitive.NilObjectID
	}

	// If event has an owner id, check if user has permissions to edit event
	if event.OwnerId != primitive.NilObjectID {
		if event.OwnerId != ownerId {
			c.JSON(http.StatusForbidden, responses.Error{Error: errs.UserNotEventOwner})
			return
		}
	}

	// Update event
	event.Name = payload.Name
	event.Duration = payload.Duration
	event.Dates = payload.Dates
	event.NotificationsEnabled = payload.NotificationsEnabled
	event.Type = payload.Type

	// Update remindees
	if event.Type == models.DOW || event.Type == models.SPECIFIC_DATES {
		origRemindees := utils.Coalesce(event.Remindees)
		updatedRemindees := make([]models.Remindee, 0)
		added, removed, kept := utils.FindAddedRemovedKept(payload.Remindees, utils.Map(origRemindees, func(r models.Remindee) string { return r.Email }))

		// Determine owner name
		var ownerName string
		if event.OwnerId == primitive.NilObjectID {
			ownerName = "Somebody"
		} else {
			owner := db.GetUserById(event.OwnerId.Hex())
			ownerName = owner.FirstName
		}

		for _, keptEmail := range kept {
			updatedRemindees = append(updatedRemindees, origRemindees[keptEmail.Index])
		}

		for _, addedEmail := range added {
			// Schedule email tasks
			taskIds := gcloud.CreateEmailTask(addedEmail.Value, ownerName, event.Name, event.GetId())
			updatedRemindees = append(updatedRemindees, models.Remindee{
				Email:     addedEmail.Value,
				TaskIds:   taskIds,
				Responded: utils.FalsePtr(),
			})
		}

		for _, removedEmail := range removed {
			// Delete email tasks
			for _, taskId := range origRemindees[removedEmail.Index].TaskIds {
				gcloud.DeleteEmailTask(taskId)
			}
		}

		event.Remindees = &updatedRemindees
	}

	// Update attendees
	if event.Type == models.GROUP {
		origAttendees := utils.Coalesce(event.Attendees)
		updatedAttendees := make([]models.Attendee, 0)
		added, _, kept := utils.FindAddedRemovedKept(payload.Attendees, utils.Map(origAttendees, func(a models.Attendee) string { return a.Email }))

		// Determine owner name
		var ownerName string
		if event.OwnerId != primitive.NilObjectID {
			owner := db.GetUserById(event.OwnerId.Hex())
			ownerName = owner.FirstName

			// Keep owner in attendees array
			if ownerIndex := utils.Find(origAttendees, func(a models.Attendee) bool { return a.Email == owner.Email }); ownerIndex != -1 {
				updatedAttendees = append(updatedAttendees, origAttendees[ownerIndex])
			}
		} else {
			ownerName = "Somebody"
		}

		for _, keptEmail := range kept {
			updatedAttendees = append(updatedAttendees, origAttendees[keptEmail.Index])
		}

		for _, addedEmail := range added {
			// Send invite email
			availabilityGroupInviteEmailId := 9
			listmonk.SendEmailAddSubscriberIfNotExist(addedEmail.Value, availabilityGroupInviteEmailId, bson.M{
				"ownerName": ownerName,
				"groupName": event.Name,
				"groupUrl":  fmt.Sprintf("%s/g/%s", utils.GetBaseUrl(), event.GetId()),
			})
			updatedAttendees = append(updatedAttendees, models.Attendee{
				Email:    addedEmail.Value,
				Declined: utils.FalsePtr(),
			})
		}

		event.Attendees = &updatedAttendees
	}

	// Update event object
	_, err := db.EventsCollection.UpdateOne(
		context.Background(),
		bson.M{
			"_id": event.Id,
		},
		bson.M{
			"$set": event,
		},
	)

	if err != nil {
		logger.StdErr.Panicln(err)
	}

	c.Status(http.StatusOK)
}

// @Summary Gets an event based on its id
// @Tags events
// @Produce json
// @Param eventId path string true "Event ID"
// @Success 200 {object} models.Event
// @Router /events/{eventId} [get]
func getEvent(c *gin.Context) {
	eventId := c.Param("eventId")
	event := db.GetEventByEitherId(eventId)

	if event == nil {
		c.JSON(http.StatusNotFound, responses.Error{Error: errs.EventNotFound})
		return
	}

	// Populate user fields
	for userId, response := range event.Responses {
		user := db.GetUserById(userId)
		if user == nil {
			if len(response.Name) == 0 {
				// User was deleted
				delete(event.Responses, userId)
				continue
			} else {
				// User is guest
				userId = response.Name
				response.User = &models.User{
					FirstName: response.Name,
				}
			}
		} else {
			response.User = user
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
// @Param payload body object{availability=[]string,guest=bool,name=string,useCalendarAvailability=bool,enabledCalendars=map[string][]string} true "Object containing info about the event response to update"
// @Success 200
// @Router /events/{eventId}/response [post]
func updateEventResponse(c *gin.Context) {
	payload := struct {
		Availability []primitive.DateTime `json:"availability"`
		Guest        *bool                `json:"guest" binding:"required"`
		Name         string               `json:"name"`

		// Calendar availability variables for Availability Groups feature
		UseCalendarAvailability *bool                `json:"useCalendarAvailability"`
		EnabledCalendars        *map[string][]string `json:"enabledCalendars"`
	}{}
	if err := c.Bind(&payload); err != nil {
		return
	}
	session := sessions.Default(c)
	eventId := c.Param("eventId")
	event := db.GetEventByEitherId(eventId)
	if event == nil {
		c.JSON(http.StatusNotFound, responses.Error{Error: errs.EventNotFound})
		return
	}

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
			UserId:                  utils.StringToObjectID(userIdString),
			Availability:            payload.Availability,
			UseCalendarAvailability: payload.UseCalendarAvailability,
			EnabledCalendars:        payload.EnabledCalendars,
		}

		// If event is group, set declined to false (in case user declined group in the past)
		if event.Type == models.GROUP {
			user := db.GetUserById(userIdString)
			if user != nil {
				for i, attendee := range utils.Coalesce(event.Attendees) {
					if attendee.Email == user.Email {
						(*event.Attendees)[i].Declined = utils.FalsePtr()
						break
					}
				}
			}
		}
	}

	// Check if user has responded to event before (edit response) or not (new response)
	_, userHasResponded := event.Responses[userIdString]

	// Update responses in mongodb
	event.Responses[userIdString] = &response
	_, err := db.EventsCollection.UpdateByID(
		context.Background(),
		event.Id,
		bson.M{"$set": event},
	)
	if err != nil {
		logger.StdErr.Panicln(err)
	}

	// Send email to creator of event if creator enabled it
	if utils.Coalesce(event.NotificationsEnabled) && !userHasResponded && userIdString != event.OwnerId.Hex() {
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
			someoneRespondedEmailId := 10
			listmonk.SendEmail(creator.Email, someoneRespondedEmailId, bson.M{
				"eventName":      event.Name,
				"ownerName":      creator.FirstName,
				"respondentName": respondentName,
				"eventUrl":       fmt.Sprintf("%s/e/%s", utils.GetBaseUrl(), event.GetId()),
			})
		}()
	}

	c.JSON(http.StatusOK, gin.H{})
}

// @Summary Delete the current user's availability
// @Tags events
// @Accept json
// @Produce json
// @Param eventId path string true "Event ID"
// @Param payload body object{userId=string,guest=bool,name=string} true "Object containing info about the event response to delete"
// @Success 200
// @Router /events/{eventId}/response [delete]
func deleteEventResponse(c *gin.Context) {
	payload := struct {
		UserId string `json:"userId"`
		Guest  *bool  `json:"guest" binding:"required"`
		Name   string `json:"name"`
	}{}
	if err := c.Bind(&payload); err != nil {
		return
	}
	session := sessions.Default(c)
	eventId := c.Param("eventId")
	event := db.GetEventByEitherId(eventId)
	if event == nil {
		c.JSON(http.StatusNotFound, responses.Error{Error: errs.EventNotFound})
		return
	}

	if *payload.Guest {
		// Delete response
		delete(event.Responses, payload.Name)
	} else {
		userIdInterface := session.Get("userId")
		if userIdInterface == nil {
			c.JSON(http.StatusUnauthorized, responses.Error{Error: errs.NotSignedIn})
			c.Abort()
			return
		}
		userIdString := userIdInterface.(string)

		// Don't allow user to delete availability of other users if they aren't the owner of the event
		if payload.UserId != userIdString && event.OwnerId.Hex() != userIdString {
			c.JSON(http.StatusForbidden, responses.Error{Error: errs.UserNotEventOwner})
			c.Abort()
			return
		}

		// Delete response
		delete(event.Responses, payload.UserId)

		// If this event is a Group, also make the attendee "leave the group" by setting "declined" to true
		if event.Type == models.GROUP {
			user := db.GetUserById(userIdString)
			if user != nil {
				for i, attendee := range utils.Coalesce(event.Attendees) {
					if attendee.Email == user.Email {
						(*event.Attendees)[i].Declined = utils.TruePtr()
						break
					}
				}
			}
		}
	}

	// Update responses in mongodb
	_, err := db.EventsCollection.UpdateByID(
		context.Background(),
		event.Id,
		bson.M{
			"$set": event,
		},
	)
	if err != nil {
		logger.StdErr.Panicln(err)
	}

	c.JSON(http.StatusOK, gin.H{})
}

// @Summary Mark the user as having responded to this event
// @Tags events
// @Accept json
// @Produce json
// @Param eventId path string true "Event ID"
// @Param payload body object{email=string} true "Object containing the user's email"
// @Success 200
// @Router /events/{eventId}/responded [post]
func userResponded(c *gin.Context) {
	payload := struct {
		Email string `json:"email" binding:"required"`
	}{}
	if err := c.Bind(&payload); err != nil {
		return
	}

	// Fetch event
	eventId := c.Param("eventId")
	event := db.GetEventByEitherId(eventId)
	if event == nil {
		c.JSON(http.StatusNotFound, responses.Error{Error: errs.EventNotFound})
		return
	}

	// Update responded boolean for the given email
	if event.Remindees == nil {
		c.JSON(http.StatusNotFound, responses.Error{Error: errs.RemindeeEmailNotFound})
		return
	}
	index := utils.Find(*event.Remindees, func(r models.Remindee) bool {
		return r.Email == payload.Email
	})
	if index == -1 {
		c.JSON(http.StatusNotFound, responses.Error{Error: errs.RemindeeEmailNotFound})
		return
	}
	if *(*event.Remindees)[index].Responded {
		// If remindee has already responded, just return and don't update db
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	(*event.Remindees)[index].Responded = utils.TruePtr()

	// Delete the reminder email tasks
	for _, taskId := range (*event.Remindees)[index].TaskIds {
		gcloud.DeleteEmailTask(taskId)
	}

	// Update event in database
	db.EventsCollection.UpdateByID(context.Background(), event.Id, bson.M{
		"$set": event,
	})

	// Email owner of event if all remindees have responded
	everyoneResponded := true
	for _, remindee := range *event.Remindees {
		if !*remindee.Responded {
			everyoneResponded = false
			break
		}
	}
	if everyoneResponded {
		// Get owner
		owner := db.GetUserById(event.OwnerId.Hex())

		// Get event url
		var baseUrl string
		if utils.IsRelease() {
			baseUrl = "https://schej.it"
		} else {
			baseUrl = "http://localhost:8080"
		}
		eventUrl := fmt.Sprintf("%s/e/%s", baseUrl, eventId)

		// Send email
		everyoneRespondedEmailTemplateId := 8
		listmonk.SendEmail(owner.Email, everyoneRespondedEmailTemplateId, bson.M{
			"eventName": event.Name,
			"eventUrl":  eventUrl,
		})
	}

	c.JSON(http.StatusOK, gin.H{})
}

// @Summary Decline the current user's invite to the event
// @Tags events
// @Accept json
// @Produce json
// @Param eventId path string true "Event ID"
// @Success 200
// @Router /events/{eventId}/decline [post]
func declineInvite(c *gin.Context) {
	// Fetch event
	eventId := c.Param("eventId")
	event := db.GetEventById(eventId)
	if event == nil {
		c.JSON(http.StatusNotFound, responses.Error{Error: errs.EventNotFound})
		return
	}

	// Ensure that event is a group
	if event.Type != models.GROUP {
		c.JSON(http.StatusBadRequest, responses.Error{Error: errs.EventNotGroup})
		return
	}

	// Get current user
	userInterface, _ := c.Get("authUser")
	user := userInterface.(*models.User)

	// Check if user is in attendees array
	index := utils.Find(utils.Coalesce(event.Attendees), func(a models.Attendee) bool {
		return user.Email == a.Email
	})
	if index == -1 {
		// User not in attendees array
		c.JSON(http.StatusNotFound, responses.Error{Error: errs.AttendeeEmailNotFound})
		return
	}

	// Decline invite
	(*event.Attendees)[index].Declined = utils.TruePtr()

	// Update event in database
	_, err := db.EventsCollection.UpdateByID(context.Background(), event.Id, bson.M{
		"$set": event,
	})
	if err != nil {
		logger.StdErr.Panicln(err)
	}

	c.JSON(http.StatusOK, gin.H{})
}

// @Summary Return a map mapping user id to their calendar events that they have enabled for the given time range
// @Tags events
// @Accept json
// @Produce json
// @Param eventId path string true "Event ID"
// @Param timeMin query string true "Lower bound for event's start time to filter by"
// @Param timeMax query string true "Upper bound for event's end time to filter by"
// @Success 200 {object} map[string]map[string]calendar.CalendarEventsWithError
// @Router /events/{eventId}/calendar-availabilities [get]
func getCalendarAvailabilities(c *gin.Context) {
	// Bind query parameters
	payload := struct {
		TimeMin time.Time `form:"timeMin" binding:"required"`
		TimeMax time.Time `form:"timeMax" binding:"required"`
	}{}
	if err := c.Bind(&payload); err != nil {
		return
	}

	// Fetch event
	eventId := c.Param("eventId")
	event := db.GetEventById(eventId)
	if event == nil {
		c.JSON(http.StatusNotFound, responses.Error{Error: errs.EventNotFound})
		return
	}

	// Ensure that event is a group
	if event.Type != models.GROUP {
		c.JSON(http.StatusBadRequest, responses.Error{Error: errs.EventNotGroup})
		return
	}

	// Fetch calendar events for each of the attendees
	numCalendarEventsRequests := 0
	calendarEventsChan := make(chan struct {
		UserId string
		Events map[string]calendar.CalendarEventsWithError
	})
	for userId, response := range event.Responses {
		if utils.Coalesce(response.UseCalendarAvailability) {
			user := db.GetUserById(userId)
			if user != nil {
				numCalendarEventsRequests++

				// Construct enabled accounts set
				enabledAccounts := make([]string, 0)
				for email := range utils.Coalesce(response.EnabledCalendars) {
					enabledAccounts = append(enabledAccounts, email)
				}

				// Fetch calendar events
				go func() {
					calendarEvents, _ := calendar.GetUsersCalendarEvents(user, utils.ArrayToSet(enabledAccounts), payload.TimeMin, payload.TimeMax)
					calendarEventsChan <- struct {
						UserId string
						Events map[string]calendar.CalendarEventsWithError
					}{
						UserId: user.Id.Hex(),
						Events: calendarEvents,
					}
				}()
			}
		}
	}

	// Create a map mapping user id to the calendar events of that user
	userIdToCalendarEvents := make(map[string][]models.CalendarEvent)
	for i := 0; i < numCalendarEventsRequests; i++ {
		calendarEvents := <-calendarEventsChan
		userIdToCalendarEvents[calendarEvents.UserId] = make([]models.CalendarEvent, 0)
		for _, events := range calendarEvents.Events {
			userIdToCalendarEvents[calendarEvents.UserId] = append(userIdToCalendarEvents[calendarEvents.UserId], events.CalendarEvents...)
		}
	}

	// Filter and format calendar events
	authUser := utils.GetAuthUser(c)
	for userId, calendarEvents := range userIdToCalendarEvents {
		// Construct enabled calendar ids set
		enabledCalendarIdsArr := make([]string, 0)
		for _, calendarIds := range utils.Coalesce(event.Responses[userId].EnabledCalendars) {
			enabledCalendarIdsArr = append(enabledCalendarIdsArr, calendarIds...)
		}
		enabledCalendarIds := utils.ArrayToSet(enabledCalendarIdsArr)

		// Update calendar events
		updatedCalendarEvents := make([]models.CalendarEvent, 0)
		for _, calendarEvent := range calendarEvents {
			// Get rid of events on sub calendars that aren't enabled
			if _, ok := enabledCalendarIds[calendarEvent.CalendarId]; !ok {
				continue
			}

			// Redact event names of other users
			if authUser.Id.Hex() != userId {
				calendarEvent.Summary = "BUSY"
			}

			updatedCalendarEvents = append(updatedCalendarEvents, calendarEvent)
		}
		userIdToCalendarEvents[userId] = updatedCalendarEvents
	}

	c.JSON(http.StatusOK, userIdToCalendarEvents)
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

// @Summary Duplicate event
// @Tags events
// @Produce json
// @Param eventId path string true "Event ID"
// @Param payload body object{eventName=string,copyAvailability=bool} true "Object containing options for the duplicated event"
// @Success 200
// @Router /events/{eventId}/duplicate [post]
func duplicateEvent(c *gin.Context) {
	payload := struct {
		EventName        string `json:"eventName" binding:"required"`
		CopyAvailability *bool  `json:"copyAvailability" binding:"required"`
	}{}
	if err := c.Bind(&payload); err != nil {
		return
	}

	eventId := c.Param("eventId")
	userInterface, _ := c.Get("authUser")
	user := userInterface.(*models.User)

	// Get event
	event := db.GetEventByEitherId(eventId)
	if event == nil {
		c.Status(http.StatusBadRequest)
		return
	}

	// Make sure user has permission to duplicate this event
	if event.OwnerId != user.Id {
		c.Status(http.StatusForbidden)
		return
	}

	// Update event
	event.Id = primitive.NewObjectID()
	event.Name = payload.EventName
	if !*payload.CopyAvailability {
		event.Responses = make(map[string]*models.Response)
	}

	// Generate short id
	shortId := db.GenerateShortEventId(event.Id)
	event.ShortId = &shortId

	// Insert new event
	result, err := db.EventsCollection.InsertOne(context.Background(), event)
	if err != nil {
		logger.StdErr.Panicln(err)
	}

	insertedId := result.InsertedID.(primitive.ObjectID).Hex()
	c.JSON(http.StatusCreated, gin.H{"eventId": insertedId, "shortId": shortId})
}
