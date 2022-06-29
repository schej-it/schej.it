/* The /friends group contains all the routes related to adding/removing/requesting friends */
package routes

import (
	"context"
	"net/http"
	"time"

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
)

func InitFriends(router *gin.Engine) {
	friendsRouter := router.Group("/friends")
	friendsRouter.Use(middleware.AuthRequired())

	friendsRouter.GET("")
	friendsRouter.DELETE("/:id")
	friendsRouter.GET("/requests")
	friendsRouter.POST("/requests", createFriendRequest)
	friendsRouter.POST("/requests/:id/accept", acceptFriendRequest)
	friendsRouter.POST("/requests/:id/reject", rejectFriendRequest)
	friendsRouter.DELETE("/requests/:id", deleteFriendRequest)
}

// @Summary Creates a new friend request
// @Tags friends
// @Accept json
// @Produce json
// @Param from body string true "The sender of the friend request"
// @Param to body string true "The recipient of the friend request"
// @Success 201 {object} models.FriendRequest
// @Router /friends/requests [post]
func createFriendRequest(c *gin.Context) {
	payload := struct {
		From primitive.ObjectID `json:"from" binding:"required"`
		To   primitive.ObjectID `json:"to" binding:"required"`
	}{}
	if err := c.Bind(&payload); err != nil {
		return
	}

	// Check if user is allowed to create this friend request
	userInterface, _ := c.Get("authUser")
	user := userInterface.(*models.User)
	if user.Id != payload.From {
		c.JSON(http.StatusForbidden, gin.H{})
		return
	}

	if result := db.FriendRequestsCollection.FindOne(context.Background(), bson.M{
		"from": payload.To,
		"to":   payload.From,
	}); result.Err() != mongo.ErrNoDocuments {
		// Friend request already exists, accept it
		var friendRequest models.FriendRequest
		result.Decode(&friendRequest)
		_acceptFriendRequest(c, &friendRequest)
		return
	}

	// Insert friend request
	friendRequest := models.FriendRequest{
		From:      payload.From,
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
// @Router /friends/requests/:id/accept [post]
func acceptFriendRequest(c *gin.Context) {
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

	// Update friend arrays of both From and To user
	db.UsersCollection.UpdateOne(context.Background(),
		bson.M{"$and": bson.A{
			bson.M{"_id": friendRequest.To},
			bson.M{"friendIds": bson.M{"$ne": bson.A{friendRequest.From}}},
		}},
		bson.M{"$push": bson.M{"friendIds": friendRequest.From}},
	)

	db.UsersCollection.UpdateOne(context.Background(),
		bson.M{"$and": bson.A{
			bson.M{"_id": friendRequest.From},
			bson.M{"friendIds": bson.M{"$ne": bson.A{friendRequest.To}}},
		}},
		bson.M{"$push": bson.M{"friendIds": friendRequest.To}},
	)

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
// @Router /friends/requests/:id/reject [post]
func rejectFriendRequest(c *gin.Context) {
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
// @Router /friends/requests/:id [delete]
func deleteFriendRequest(c *gin.Context) {

}
