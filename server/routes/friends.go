/* The /friends group contains all the routes related to adding/removing/requesting friends */
package routes

import "github.com/gin-gonic/gin"

func InitFriends(router *gin.Engine) {
	friendsRouter := router.Group("/friends")
	friendsRouter.GET("")
	friendsRouter.DELETE("/:id")
	friendsRouter.GET("/requests")
	friendsRouter.POST("/requests")
	friendsRouter.POST("/requests/:id/reject")
	friendsRouter.POST("/requests/:id/accept")
	friendsRouter.DELETE("/requests/:id")
}
