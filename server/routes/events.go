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
	"schej.it/server/errs"
	"schej.it/server/logger"
	"schej.it/server/middleware"
	"schej.it/server/models"
	"schej.it/server/responses"
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
	// eventRouter.POST("/:eventId/attendee", middleware.AuthRequired(), addAttendee)
	// eventRouter.DELETE("/:eventId/attendee", middleware.AuthRequired(), removeAttendee)
	eventRouter.DELETE("/:eventId", middleware.AuthRequired(), deleteEvent)
	eventRouter.POST("/:eventId/duplicate", middleware.AuthRequired(), duplicateEvent)
}

// @Summary Creates a new event
// @Tags events
// @Accept json
// @Produce json
// @Param payload body object{name=string,duration=float32,dates=[]string,notificationsEnabled=bool,remindees=[]string,type=models.EventType} true "Object containing info about the event to create"
// @Success 201 {object} object{eventId=string}
// @Router /events [post]
func createEvent(c *gin.Context) {
	payload := struct {
		Name                 string               `json:"name" binding:"required"`
		Duration             *float32             `json:"duration" binding:"required"`
		Dates                []primitive.DateTime `json:"dates" binding:"required"`
		NotificationsEnabled *bool                `json:"notificationsEnabled" binding:"required"`
		Remindees            []string             `json:"remindees" binding:"required"`
		Type                 models.EventType     `json:"type" binding:"required"`
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

	event := models.Event{
		Id:                   primitive.NewObjectID(),
		OwnerId:              ownerId,
		Name:                 payload.Name,
		Duration:             payload.Duration,
		Dates:                payload.Dates,
		NotificationsEnabled: *payload.NotificationsEnabled,
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
			taskIds := gcloud.CreateEmailTask(email, ownerName, payload.Name, event.Id.Hex())
			remindees = append(remindees, models.Remindee{
				Email:     email,
				TaskIds:   taskIds,
				Responded: utils.FalsePtr(),
			})
		}

		event.Remindees = remindees
	}

	result, err := db.EventsCollection.InsertOne(context.Background(), event)
	if err != nil {
		logger.StdErr.Panicln(err)
	}

	insertedId := result.InsertedID.(primitive.ObjectID).Hex()

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
// @Param payload body object{name=string,duration=float32,dates=[]string,notificationsEnabled=bool,remindees=[]string,type=models.EventType} true "Object containing info about the event to update"
// @Success 200
// @Router /events/{eventId} [put]
func editEvent(c *gin.Context) {
	payload := struct {
		Name                 string               `json:"name" binding:"required"`
		Duration             *float32             `json:"duration" binding:"required"`
		Dates                []primitive.DateTime `json:"dates" binding:"required"`
		NotificationsEnabled *bool                `json:"notificationsEnabled" binding:"required"`
		Remindees            []string             `json:"remindees" binding:"required"`
		Type                 models.EventType     `json:"type" binding:"required"`
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
	event := db.GetEventByEitherId(eventId)

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
	event.NotificationsEnabled = *payload.NotificationsEnabled
	event.Type = payload.Type

	//
	// Update remindees
	//
	updatedRemindees := make([]models.Remindee, 0)

	// Find remindees to delete
	for i, remindee := range event.Remindees {
		found := false
		for _, email := range payload.Remindees {
			if email == remindee.Email {
				found = true
				break
			}
		}

		if found {
			// If remindee is found in updated remindees array, add to updatedRemindees
			updatedRemindees = append(updatedRemindees, event.Remindees[i])
		} else {
			// If remindee not found in the updated remindees array, delete remindee

			// Delete email tasks
			for _, taskId := range event.Remindees[i].TaskIds {
				gcloud.DeleteEmailTask(taskId)
			}
		}
	}

	// Find remindees to add
	for _, email := range payload.Remindees {
		found := false
		for _, remindee := range updatedRemindees {
			if email == remindee.Email {
				found = true
				break
			}
		}

		if !found {
			// Update remindees array contains a new remindee, so add a new remindee

			// Determine owner name
			var ownerName string
			if event.OwnerId != primitive.NilObjectID {
				owner := db.GetUserById(event.OwnerId.Hex())
				ownerName = owner.FirstName
			} else {
				ownerName = "Somebody"
			}

			// Schedule email tasks
			taskIds := gcloud.CreateEmailTask(email, ownerName, event.Name, event.Id.Hex())
			updatedRemindees = append(updatedRemindees, models.Remindee{
				Email:     email,
				TaskIds:   taskIds,
				Responded: utils.FalsePtr(),
			})
		}
	}

	event.Remindees = updatedRemindees

	// Update event object
	_, err = db.EventsCollection.UpdateOne(
		context.Background(),
		bson.M{
			"_id": objectId,
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
			userId = response.Name
			response.User = &models.User{
				FirstName: response.Name,
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
// @Param payload body object{availability=[]string,guest=bool,name=string,attendeeEmail=string} true "Object containing info about the event response to update"
// @Success 200
// @Router /events/{eventId}/response [post]
func updateEventResponse(c *gin.Context) {
	payload := struct {
		Availability []primitive.DateTime `json:"availability" binding:"required"`
		Guest        *bool                `json:"guest" binding:"required"`
		Name         string               `json:"name"`
	}{}
	if err := c.Bind(&payload); err != nil {
		return
	}
	session := sessions.Default(c)
	eventId := c.Param("eventId")
	event := db.GetEventByEitherId(eventId)

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
		bson.A{
			utils.UpdateEventResponseAggregation(userIdString, response),
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

	var userToDelete string
	if *payload.Guest {
		userToDelete = payload.Name
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

		userToDelete = payload.UserId
	}

	// Update responses in mongodb
	_, err := db.EventsCollection.UpdateByID(
		context.Background(),
		utils.StringToObjectID(eventId),
		bson.A{
			utils.DeleteEventResponseAggregation(userToDelete),
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

	// Update responded boolean for the given email
	index := utils.Find(event.Remindees, func(r models.Remindee) bool {
		return r.Email == payload.Email
	})
	if index == -1 {
		c.JSON(http.StatusNotFound, responses.Error{Error: errs.RemindeeEmailNotFound})
		return
	}
	if *event.Remindees[index].Responded {
		// If remindee has already responded, just return and don't update db
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	event.Remindees[index].Responded = utils.TruePtr()

	// Delete the reminder email tasks
	for _, taskId := range event.Remindees[index].TaskIds {
		gcloud.DeleteEmailTask(taskId)
	}

	// Update event in database
	db.EventsCollection.UpdateByID(context.Background(), event.Id, bson.M{
		"$set": event,
	})

	// Email owner of event if all remindees have responded
	everyoneResponded := true
	for _, remindee := range event.Remindees {
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

// @Summary [UNUSED] Adds an attendee to the event's list of attendees
// @Tags events
// @Accept json
// @Produce json
// @Param eventId path string true "Event ID"
// @Param payload body object{email=string} true "Object containing info about the attendee to add"
// @Success 200
// @Router /events/{eventId}/attendee [post]
// func addAttendee(c *gin.Context) {
// 	// UNUSED
// 	payload := struct {
// 		Email string `json:"email" binding:"required"`
// 	}{}
// 	if err := c.Bind(&payload); err != nil {
// 		return
// 	}
// 	session := sessions.Default(c)
// 	eventIdString := c.Param("eventId")
// 	eventId, err := primitive.ObjectIDFromHex(eventIdString)
// 	if err != nil {
// 		logger.StdErr.Panicln(err)
// 	}

// 	userIdInterface := session.Get("userId")
// 	userIdString := userIdInterface.(string)
// 	userId, err := primitive.ObjectIDFromHex(userIdString)
// 	if err != nil {
// 		logger.StdErr.Panicln(err)
// 	}

// 	_, err = db.EventsCollection.UpdateOne(context.Background(), bson.M{
// 		"_id":     eventId,
// 		"ownerId": userId,
// 	}, bson.M{
// 		"$addToSet": bson.M{"attendees": payload.Email},
// 	})
// 	if err != nil {
// 		logger.StdErr.Panicln(err)
// 	}

// 	c.JSON(http.StatusOK, gin.H{})
// }

// @Summary [UNUSED] Removes an attendee from the event's list of attendees
// @Tags events
// @Accept json
// @Produce json
// @Param eventId path string true "Event ID"
// @Param payload body object{email=string} true "Object containing info about the attendee to remove"
// @Success 200
// @Router /events/{eventId}/attendee [delete]
// func removeAttendee(c *gin.Context) {
// 	// UNUSED
// 	payload := struct {
// 		Email string `json:"email" binding:"required"`
// 	}{}
// 	if err := c.Bind(&payload); err != nil {
// 		return
// 	}
// 	session := sessions.Default(c)
// 	eventIdString := c.Param("eventId")
// 	eventId, err := primitive.ObjectIDFromHex(eventIdString)
// 	if err != nil {
// 		logger.StdErr.Panicln(err)
// 	}

// 	userIdInterface := session.Get("userId")
// 	userIdString := userIdInterface.(string)
// 	userId, err := primitive.ObjectIDFromHex(userIdString)
// 	if err != nil {
// 		logger.StdErr.Panicln(err)
// 	}

// 	_, err = db.EventsCollection.UpdateOne(context.Background(), bson.M{
// 		"_id":     eventId,
// 		"ownerId": userId,
// 	}, bson.M{
// 		"$pull": bson.M{"attendees": payload.Email},
// 	})
// 	if err != nil {
// 		logger.StdErr.Panicln(err)
// 	}

// 	c.JSON(http.StatusOK, gin.H{})
// }
