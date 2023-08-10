/* The /auth group contains all the routes to sign in and sign out */
package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"schej.it/server/db"
	"schej.it/server/logger"
	"schej.it/server/middleware"
	"schej.it/server/models"
	"schej.it/server/slackbot"
	"schej.it/server/utils"
)

func InitAuth(router *gin.Engine) {
	authRouter := router.Group("/auth")

	authRouter.POST("/sign-in", signIn)
	authRouter.POST("/sign-in-mobile", signInMobile)
	authRouter.POST("/sign-out", signOut)
	authRouter.GET("/status", middleware.AuthRequired(), getStatus)
}

// @Summary Signs user in
// @Description Signs user in and sets the access token session variable
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body object{code=string,timezoneOffset=int} true "Object containing the Google authorization code and the user's timezone offset"
// @Success 200
// @Router /auth/sign-in [post]
func signIn(c *gin.Context) {
	payload := struct {
		Code           string `json:"code" binding:"required"`
		TimezoneOffset int    `json:"timezoneOffset" binding:"required"`
	}{}
	if err := c.BindJSON(&payload); err != nil {
		return
	}

	// Call Google oauth token endpoint
	var redirectUri string
	if utils.IsRelease() {
		redirectUri = "https://schej.it/auth"
	} else {
		redirectUri = "http://localhost:8080/auth"
	}
	values := url.Values{
		"client_id":     {os.Getenv("CLIENT_ID")},
		"client_secret": {os.Getenv("CLIENT_SECRET")},
		"code":          {payload.Code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {redirectUri},
	}
	resp, err := http.PostForm(
		"https://oauth2.googleapis.com/token",
		values,
	)
	if err != nil {
		logger.StdErr.Panicln(err)
	}
	defer resp.Body.Close()

	res := struct {
		AccessToken      string `json:"access_token"`
		IdToken          string `json:"id_token"`
		ExpiresIn        int    `json:"expires_in"`
		RefreshToken     string `json:"refresh_token"`
		Scope            string `json:"scope"`
		TokenType        string `json:"token_type"`
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
	}{}

	json.NewDecoder(resp.Body).Decode(&res)
	if len(res.Error) > 0 {
		data, _ := json.MarshalIndent(res, "", "  ")
		logger.StdErr.Panicln(string(data))
	}

	signInHelper(c, res.AccessToken, res.IdToken, res.ExpiresIn, res.RefreshToken, payload.TimezoneOffset, models.WEB)

	c.JSON(http.StatusOK, gin.H{})
}

// @Summary Signs user in from mobile
// @Description Signs user in and sets the access token session variable
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body object{timezoneOffset=int,accessToken=string,idToken=string,expiresIn=int,refreshToken=string,scope=string} true "Object containing the Google authorization code and the user's timezone offset"
// @Success 200
// @Router /auth/sign-in-mobile [post]
func signInMobile(c *gin.Context) {
	payload := struct {
		TimezoneOffset int                    `json:"timezoneOffset" binding:"required"`
		AccessToken    string                 `json:"accessToken" binding:"required"`
		IdToken        string                 `json:"idToken" binding:"required"`
		ExpiresIn      int                    `json:"expiresIn" binding:"required"`
		RefreshToken   string                 `json:"refreshToken" binding:"required"`
		TokenOrigin    models.TokenOriginType `json:"tokenOrigin" binding:"required"`
	}{}
	if err := c.BindJSON(&payload); err != nil {
		return
	}

	signInHelper(c, payload.AccessToken, payload.IdToken, payload.ExpiresIn, payload.RefreshToken, payload.TimezoneOffset, payload.TokenOrigin)

	c.JSON(http.StatusOK, gin.H{})
}

// Helper function to sign user in with the given parameters from the google oauth route
func signInHelper(c *gin.Context, accessToken string, idToken string, expiresIn int, refreshToken string, timezoneOffset int, tokenOrigin models.TokenOriginType) {
	// Get access token expire time
	accessTokenExpireDate := utils.GetAccessTokenExpireDate(expiresIn)

	// Get user info from JWT
	claims := utils.ParseJWT(idToken)
	email, _ := claims.GetStr("email")
	firstName, _ := claims.GetStr("given_name")
	lastName, _ := claims.GetStr("family_name")
	picture, _ := claims.GetStr("picture")

	// Create user object to create new user or update existing user
	userData := models.User{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Picture:   picture,

		// TODO: make sure this doesn't result in bugs, make sure doesn't reset to 0 randomly
		Visibility: 0,

		// FriendIds: make([]primitive.ObjectID, 0),
		// Calendars: make(map[string]models.Calendar),

		AccessToken:           accessToken,
		AccessTokenExpireDate: primitive.NewDateTimeFromTime(accessTokenExpireDate),
		RefreshToken:          refreshToken,

		TimezoneOffset: timezoneOffset,
		TokenOrigin:    tokenOrigin,
	}

	// Update user if exists
	updateResult := db.UsersCollection.FindOneAndUpdate(
		context.Background(),
		bson.M{"email": email},
		bson.M{"$set": userData},
	)
	var userId primitive.ObjectID
	if updateResult.Err() == mongo.ErrNoDocuments {
		// User doesn't exist, create a new user
		res, err := db.UsersCollection.InsertOne(context.Background(), userData)
		if err != nil {
			logger.StdErr.Panicln(err)
		}

		userId = res.InsertedID.(primitive.ObjectID)

		slackbot.SendTextMessage(fmt.Sprintf(":wave: %s %s (%s) has joined schej.it!", firstName, lastName, email))
		utils.AddUserToMailjet(email, firstName, lastName, picture)
	} else {
		// User does exist, get user id
		var user models.User
		if err := updateResult.Decode(&user); err != nil {
			logger.StdErr.Panicln(err)
		}

		userId = user.Id
	}

	// Set session variables
	session := sessions.Default(c)
	session.Set("userId", userId.Hex())
	session.Save()
}

// @Summary Signs user out
// @Description Signs user out and deletes the session
// @Tags auth
// @Accept json
// @Produce json
// @Success 200
// @Router /auth/sign-out [post]
func signOut(c *gin.Context) {
	// Delete session
	session := sessions.Default(c)
	session.Delete("userId")
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
