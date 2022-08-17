/* The /friends group contains all the routes related to adding/removing/requesting friends */
package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"schej.it/server/db"
	"schej.it/server/errs"
	"schej.it/server/logger"
	"schej.it/server/middleware"
	"schej.it/server/models"
	"schej.it/server/responses"
	"schej.it/server/utils"
)

func InitFriends(router *gin.Engine) {
	friendsRouter := router.Group("/friends")
	friendsRouter.Use(middleware.AuthRequired())

	friendsRouter.GET("", getFriends)
	friendsRouter.GET("/:id/schedule", getFriendsSchedule)
	friendsRouter.DELETE("/:id", deleteFriend)
	friendsRouter.GET("/requests", getFriendRequests)
	friendsRouter.POST("/requests", createFriendRequest)
	friendsRouter.POST("/requests/:id/accept", acceptFriendRequest)
	friendsRouter.POST("/requests/:id/reject", rejectFriendRequest)
	friendsRouter.DELETE("/requests/:id", deleteFriendRequest)
}

// @Summary Gets all of users current friends
// @Tags friends
// @Accept json
// @Produce json
// @Success 200
// @Router /friends [get]
func getFriends(c *gin.Context) {

	session := sessions.Default(c)
	userId := utils.GetUserId(session)

	var result []models.User

	// Find and populate
	cursor, err := db.UsersCollection.Aggregate(context.Background(), []bson.M{
		{"$match": bson.M{
			"_id": userId,
		}},
		{"$lookup": bson.M{
			"from":         "users",
			"localField":   "friendIds",
			"foreignField": "_id",
			"as":           "friends",
		}},
		{"$project": bson.M{
			"friends._id":        1,
			"friends.firstName":  1,
			"friends.picture":    1,
			"friends.lastName":   1,
			"friends.email":      1,
			"friends.visibility": 1,
		}},
	})
	if err != nil {
		logger.StdErr.Panicln(err)
	}
	if err := cursor.All(context.Background(), &result); err != nil {
		logger.StdErr.Panicln(err)
	}

	c.JSON(http.StatusOK, result[0].Friends)

}

// @Summary Returns the specified friend's schedule
// @Tags friends
// @Accept json
// @Produce json
// @Param id path string true "ID of friend"
// @Param timeMin query string true "Lower bound for event's start time to filter by"
// @Param timeMax query string true "Upper bound for event's end time to filter by"
// @Success 200 {object} []models.CalendarEvent
// @Router /friends/{id}/schedule [get]
func getFriendsSchedule(c *gin.Context) {
	// Bind query parameters
	payload := struct {
		TimeMin time.Time `form:"timeMin" binding:"required"`
		TimeMax time.Time `form:"timeMax" binding:"required"`
	}{}
	if err := c.Bind(&payload); err != nil {
		return
	}
	friendId := c.Param("id")

	// See if friend to check schedule of is an existing user
	friend := db.GetUserById(friendId)
	if friend == nil {
		c.JSON(http.StatusNotFound, responses.Error{Error: errs.UserDoesNotExist})
		return
	}

	userInterface, _ := c.Get("authUser")
	user := userInterface.(*models.User)
	friendObjectID, err := primitive.ObjectIDFromHex(friendId)
	if err != nil {
		logger.StdErr.Panic(err)
	}

	// Make sure user is friends with the given friend
	if !utils.Contains(user.FriendIds, friendObjectID) {
		c.JSON(http.StatusForbidden, responses.Error{Error: errs.UserNotFriends})
		return
	}

	// Get calendar events
	calendarEvents, googleErr := db.GetUsersCalendarEvents(friend, payload.TimeMin, payload.TimeMax)
	if googleErr != nil {
		c.JSON(googleErr.Code, responses.Error{Error: *googleErr})
		return
	}

	c.JSON(http.StatusOK, calendarEvents)
}

// @Summary Removes an existing friend
// @Tags friends
// @Accept json
// @Produce json
// @Param id path string true "ID of friend"
// @Success 200
// @Router /friends/{id} [delete]
func deleteFriend(c *gin.Context) {
	session := sessions.Default(c)
	userId := utils.GetUserId(session)
	friendId := c.Param("id")

	// See if friend to be deleted is an existing user
	friend := db.GetUserById(friendId)
	if friend == nil {
		c.JSON(http.StatusNotFound, responses.Error{Error: errs.UserDoesNotExist})
		return
	}

	// Remove friend from friend array
	db.UsersCollection.UpdateOne(context.Background(),
		bson.M{
			"_id": userId,
		},
		bson.M{"$pullAll": bson.M{"friendIds": friendId}},
	)

}

// @Summary Gets all the current incoming and outgoing friend requests
// @Tags friends
// @Accept json
// @Produce json
// @Success 200
// @Router /friends/requests [get]
func getFriendRequests(c *gin.Context) {
	session := sessions.Default(c)
	userId := utils.GetUserId(session)

	var result []models.FriendRequest

	// Find and populate
	cursor, err := db.FriendRequestsCollection.Aggregate(context.Background(), []bson.M{
		{"$match": bson.M{
			"$or": bson.A{
				bson.M{"to": userId},
				bson.M{"from": userId},
			},
		}},
		{"$lookup": bson.M{
			"from":         "users",
			"localField":   "to",
			"foreignField": "_id",
			"as":           "toUser",
		}},
		{"$unwind": "$toUser"},
		{"$lookup": bson.M{
			"from":         "users",
			"localField":   "from",
			"foreignField": "_id",
			"as":           "fromUser",
		}},
		{"$unwind": "$fromUser"},
	})
	if err != nil {
		logger.StdErr.Panicln(err)
	}
	if err := cursor.All(context.Background(), &result); err != nil {
		logger.StdOut.Println(result)
		logger.StdErr.Panicln(err)
	}

	c.JSON(http.StatusOK, result)
}

// @Summary Creates a new friend request
// @Tags friends
// @Accept json
// @Produce json
// @Param payload body object{to=string} true "Object specifying the user IDs of who this request is sent from and to"
// @Success 201 {object} models.FriendRequest "Friend request created"
// @Success 200 "Friend request already exists from \"to\" to \"from\", and it was accepted"
// @Router /friends/requests [post]
func createFriendRequest(c *gin.Context) {
	payload := struct {
		To primitive.ObjectID `json:"to" binding:"required"`
	}{}
	if err := c.Bind(&payload); err != nil {
		return
	}

	// Get user info
	userInterface, _ := c.Get("authUser")
	user := userInterface.(*models.User)

	// TODOS: (not essential, because frontend should prevent these errors)
	// TODO: If user is already friends with the user, return an error

	// If friend request already exists, throw an error
	if result := db.FriendRequestsCollection.FindOne(context.Background(), bson.M{
		"to":   payload.To,
		"from": user.Id,
	}); result.Err() != mongo.ErrNoDocuments {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	// If friend request already exists from the other user, accept it
	if result := db.FriendRequestsCollection.FindOne(context.Background(), bson.M{
		"from": payload.To,
		"to":   user.Id,
	}); result.Err() != mongo.ErrNoDocuments {
		var friendRequest models.FriendRequest
		result.Decode(&friendRequest)
		_acceptFriendRequest(c, &friendRequest)
		return
	}

	// Insert new friend request
	friendRequest := models.FriendRequest{
		From:      user.Id,
		To:        payload.To,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}
	result, err := db.FriendRequestsCollection.InsertOne(context.Background(), friendRequest)
	if err != nil {
		logger.StdErr.Panicln(err)
	}

	insertedId := result.InsertedID.(primitive.ObjectID)
	friendRequest.Id = insertedId

	// Populate the ToUser field
	friendRequest.ToUser = db.GetUserById(friendRequest.To.Hex()).GetProfile()

	c.JSON(http.StatusCreated, friendRequest)
}

// @Summary Accepts an existing friend request
// @Tags friends
// @Accept json
// @Produce json
// @Param id path string true "ID of the friend request"
// @Success 200
// @Router /friends/requests/{id}/accept [post]
func acceptFriendRequest(c *gin.Context) {
	// Check that the specified friend request exists
	friendRequestId := c.Param("id")
	friendRequest := db.GetFriendRequestById(friendRequestId)
	if friendRequest == nil {
		c.JSON(http.StatusNotFound, responses.Error{Error: errs.FriendRequestNotFound})
		return
	}

	_acceptFriendRequest(c, friendRequest)
}

// Helper function for friend request route
func _acceptFriendRequest(c *gin.Context, friendRequest *models.FriendRequest) {
	// Check if the "To" user id matches the current user's id
	userInterface, _ := c.Get("authUser")
	user := userInterface.(*models.User)
	if user.Id != friendRequest.To {
		c.JSON(http.StatusForbidden, gin.H{})
		return
	}

	// Update friend array of the To user
	_, err := db.UsersCollection.UpdateOne(context.Background(),
		bson.M{"_id": friendRequest.To},
		bson.M{"$addToSet": bson.M{"friendIds": friendRequest.From}},
	)
	if err != nil {
		logger.StdErr.Panicln(err)
	}

	// Update friend array of the From user
	_, err = db.UsersCollection.UpdateOne(context.Background(),
		bson.M{"_id": friendRequest.From},
		bson.M{"$addToSet": bson.M{"friendIds": friendRequest.To}},
	)
	if err != nil {
		logger.StdErr.Panicln(err)
	}

	// Delete friend request
	db.DeleteFriendRequestById(friendRequest.Id.Hex())

	c.JSON(http.StatusOK, gin.H{})
}

// @Summary Rejects an existing friend request
// @Tags friends
// @Accept json
// @Produce json
// @Param id path string true "ID of the friend request"
// @Success 200
// @Router /friends/requests/{id}/reject [post]
func rejectFriendRequest(c *gin.Context) {
	// Check that the specified friend request exists
	friendRequestId := c.Param("id")
	friendRequest := db.GetFriendRequestById(friendRequestId)
	if friendRequest == nil {
		c.JSON(http.StatusNotFound, responses.Error{Error: errs.FriendRequestNotFound})
		return
	}

	// Check if the "To" user id matches the current user's id
	userInterface, _ := c.Get("authUser")
	user := userInterface.(*models.User)
	if user.Id != friendRequest.To {
		c.JSON(http.StatusForbidden, gin.H{})
		return
	}

	// Delete friend request
	db.DeleteFriendRequestById(friendRequestId)

	c.JSON(http.StatusOK, gin.H{})
}

// @Summary Delete's a friend request created by the current user
// @Tags friends
// @Accept json
// @Produce json
// @Param id path string true "ID of the friend request"
// @Success 200
// @Router /friends/requests/{id} [delete]
func deleteFriendRequest(c *gin.Context) {
	// Check that the specified friend request exists
	friendRequestId := c.Param("id")
	friendRequest := db.GetFriendRequestById(friendRequestId)
	if friendRequest == nil {
		c.JSON(http.StatusNotFound, responses.Error{Error: errs.FriendRequestNotFound})
		return
	}

	// Check if the "From" user id matches the current user's id
	userInterface, _ := c.Get("authUser")
	user := userInterface.(*models.User)
	if user.Id != friendRequest.From {
		c.JSON(http.StatusForbidden, gin.H{})
		return
	}

	// Delete friend request
	db.DeleteFriendRequestById(friendRequestId)

	c.JSON(http.StatusOK, gin.H{})
}
