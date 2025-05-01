/* The /analytics group contains all the routes to track analytics */
package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"schej.it/server/db"
	"schej.it/server/models"
	"schej.it/server/slackbot"
)

func InitAnalytics(router *gin.RouterGroup) {
	authRouter := router.Group("/analytics")

	authRouter.POST("/scanned-poster", scannedPoster)
	authRouter.POST("/upgrade-dialog-viewed", upgradeDialogViewed)
}

// @Summary Notifies us when poster QR code has been scanned
// @Tags analytics
// @Accept json
// @Produce json
// @Param payload body object{url=string,location=models.Location} true "Object containing the location that poster was scanned from and the url that was scanned"
// @Success 200
// @Router /analytics/scanned-poster [post]
func scannedPoster(c *gin.Context) {
	payload := struct {
		Url      string           `json:"url" binding:"required"`
		Location *models.Location `json:"location"`
	}{}
	if err := c.BindJSON(&payload); err != nil {
		return
	}

	if payload.Location != nil {
		slackbot.SendTextMessage(
			fmt.Sprintf(":face_with_monocle: Poster was scanned :face_with_monocle:\n*Location:* %s, %s, %s\n*URL:* %s",
				payload.Location.City,
				payload.Location.State,
				payload.Location.CountryCode,
				payload.Url,
			),
		)
	} else {
		slackbot.SendTextMessage(
			fmt.Sprintf(":face_with_monocle: Poster was scanned :face_with_monocle:\n*URL:* %s", payload.Url),
		)
	}

	c.JSON(http.StatusOK, gin.H{})
}

// @Summary Notifies us when user has viewed the upgrade dialog
// @Tags analytics
// @Accept json
// @Produce json
// @Param payload body object{userId=string} true "Object containing the user id"
// @Success 200
// @Router /analytics/upgrade-dialog-viewed [post]
func upgradeDialogViewed(c *gin.Context) {
	payload := struct {
		UserId string `json:"userId" binding:"required"`
		Price  string `json:"price" binding:"required"`
	}{}
	if err := c.BindJSON(&payload); err != nil {
		return
	}

	user := db.GetUserById(payload.UserId)
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	slackbot.SendTextMessageWithType(
		fmt.Sprintf(":eyes: %s %s (%s) viewed the upgrade dialog (%s)", user.FirstName, user.LastName, user.Email, payload.Price),
		slackbot.MONETIZATION,
	)

	c.JSON(http.StatusOK, gin.H{})
}
