/* The /users group contains all the routes to get information about all users */
package routes

import "github.com/gin-gonic/gin"

func InitUsers(router *gin.Engine) {
	userRouter := router.Group("/users")
	userRouter.GET("", searchUsers)
}

func searchUsers(c *gin.Context) {

}
