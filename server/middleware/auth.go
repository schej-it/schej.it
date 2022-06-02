package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"schej.it/server/db"
	"schej.it/server/errs"
	"schej.it/server/responses"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if userId is set
		session := sessions.Default(c)
		if session.Get("userId") == nil {
			// User id is not set, user is not signed in!
			c.JSON(http.StatusUnauthorized, responses.Error{Error: errs.NotSignedIn})
			c.Abort()
			return
		}

		// Check if user with user id exists
		user := db.GetUserById(session.Get("userId").(string))

		if user == nil {
			c.JSON(http.StatusUnauthorized, responses.Error{Error: errs.UserDoesNotExist})
			c.Abort()
			return
		}

		c.Set("authUser", user)

		c.Next()
	}
}
