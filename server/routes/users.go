/* The /users group contains all the routes to get information about all users */
package routes

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"schej.it/server/db"
	"schej.it/server/logger"
	"schej.it/server/models"
	"schej.it/server/utils"
)

func InitUsers(router *gin.RouterGroup) {
	userRouter := router.Group("/users")
	userRouter.GET("", searchUsers)
	userRouter.GET("/:userId", getUser)
}

// @Summary Returns users that match the search query
// @Tags users
// @Produce json
// @Param query query string true "Search query matching users' names/emails"
// @Success 200 {object} []models.User "An array of user profile objects"
// @Router /users [get]
func searchUsers(c *gin.Context) {

	// Bind query parameters
	payload := struct {
		Query *string `form:"query" binding:"required"`
	}{}
	if err := c.Bind(&payload); err != nil {
		logger.StdErr.Panicln(err)
		return
	}

	query := *payload.Query
	queryTermsRegex := make([]primitive.Regex, 0)

	for _, s := range strings.Split(query, " ") {
		r := primitive.Regex{Pattern: utils.EscapeRegExp(s), Options: "i"}
		queryTermsRegex = append(queryTermsRegex, r)
	}

	users := make([]models.User, 0)

	cursor, err := db.UsersCollection.Find(context.Background(), bson.M{
		"$expr": bson.M{
			"$reduce": bson.M{
				"input":        queryTermsRegex,
				"initialValue": true,
				"in": bson.M{
					"$and": bson.A{
						"$$value",
						bson.M{
							"$regexMatch": bson.M{
								"input": bson.M{
									"$concat": bson.A{
										"$firstName", " ", "$lastName", " ", "$email",
									},
								},
								"regex": "$$this",
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		logger.StdErr.Panicln(err)
	}
	if err := cursor.All(context.Background(), &users); err != nil {
		logger.StdErr.Panicln(err)
	}

	c.JSON(http.StatusOK, users)
}

// @Summary Returns the user by their user id
// @Tags users
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} models.User "A user profile object"
// @Router /users/{userId} [get]
func getUser(c *gin.Context) {
	userId := c.Param("userId")
	user := db.GetUserById(userId)

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	c.JSON(http.StatusOK, user)
}
