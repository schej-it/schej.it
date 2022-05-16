package middleware

import (
	"context"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"schej.it/server/db"
	"schej.it/server/errors"
	"schej.it/server/models"
	"schej.it/server/utils"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if userId is set
		session := sessions.Default(c)
		if session.Get("userId") == nil {
			// User id is not set, user is not signed in!
			c.JSON(http.StatusForbidden, gin.H{"error": errors.NotSignedIn})
			c.Abort()
			return
		}

		// Check if user with user id exists
		result := db.UsersCollection.FindOne(context.Background(), bson.M{
			"_id": utils.GetUserId(session),
		})
		if result.Err() == mongo.ErrNoDocuments {
			// User does not exist!
			c.JSON(http.StatusForbidden, gin.H{"error": errors.UserDoesNotExist})
			c.Abort()
			return
		}

		// Set auth user request variable
		var user models.User
		if err := result.Decode(&user); err != nil {
			panic(err)
		}
		c.Set("authUser", user)

		c.Next()
	}
}
