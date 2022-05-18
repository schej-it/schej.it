package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"schej.it/server/db"
	"schej.it/server/middleware"
	"schej.it/server/models"
	"schej.it/server/utils"
)

func InitAuth(router *gin.Engine) {
	authRouter := router.Group("/auth")

	authRouter.POST("/sign-in", postSignIn)
	authRouter.GET("/status", middleware.AuthRequired(), getStatus)
}

// @Summary Signs user in
// @Description Signs user in and sets the access token session variable
// @Tags auth
// @Accept json
// @Produce json
// @Param code body string true "Google authorization code"
// @Success 200
// @Router /auth/sign-in [post]
func postSignIn(c *gin.Context) {
	payload := struct {
		Code string `json:"code" binding:"required"`
	}{}
	if err := c.BindJSON(&payload); err != nil {
		return
	}

	// Get access_token
	values := url.Values{
		"client_id":     {os.Getenv("CLIENT_ID")},
		"client_secret": {os.Getenv("CLIENT_SECRET")},
		"code":          {payload.Code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {"http://localhost:8080/auth"},
	}
	resp, err := http.PostForm(
		"https://oauth2.googleapis.com/token",
		values,
	)
	if err != nil {
		panic(err)
	}
	var res gin.H
	json.NewDecoder(resp.Body).Decode(&res)

	// Get user info from JWT
	claims := utils.ParseJWT(res["id_token"].(string))
	email, _ := claims.GetStr("email")

	// Create new user object if it doesn't exist, and set the userId
	result := db.UsersCollection.FindOne(context.Background(), bson.M{"email": email})

	var userId primitive.ObjectID
	if result.Err() == mongo.ErrNoDocuments {
		firstName, _ := claims.GetStr("given_name")
		lastName, _ := claims.GetStr("family_name")
		picture, _ := claims.GetStr("picture")
		newUser := models.User{
			Email:     email,
			FirstName: firstName,
			LastName:  lastName,
			Picture:   picture,
		}
		res, err := db.UsersCollection.InsertOne(context.Background(), newUser)
		if err != nil {
			panic(err)
		}

		userId = res.InsertedID.(primitive.ObjectID)
	} else {
		var user models.User
		if err := result.Decode(&user); err != nil {
			panic(err)
		}

		userId = user.Id
	}

	// Set session variables
	session := sessions.Default(c)
	session.Set("userId", userId.Hex())
	session.Set("accessToken", res["access_token"])
	session.Save()

	c.JSON(http.StatusOK, gin.H{})
}

// @Summary Gets whether the user is signed in or not
// @Description Returns a 401 error if user is not signed in, 200 if they are
// @Tags auth
// @Success 200
// @Failure 401 {object} responses.Error "Error object"
// @Router /auth/status [get]
func getStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}