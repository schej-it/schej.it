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

	authRouter.POST("/sign-in", signIn)
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
func signIn(c *gin.Context) {
	payload := struct {
		Code string `json:"code" binding:"required"`
	}{}
	if err := c.BindJSON(&payload); err != nil {
		return
	}

	// Call Google oauth token endpoint
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
	res := struct {
		AccessToken  string `json:"access_token"`
		IdToken      string `json:"id_token"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		Scope        string `json:"scope"`
		TokenType    string `json:"token_type"`
	}{}
	json.NewDecoder(resp.Body).Decode(&res)

	// Get access token expire time
	accessTokenExpireDate := utils.GetAccessTokenExpireDate(res.ExpiresIn)

	// Get user info from JWT
	claims := utils.ParseJWT(res.IdToken)
	email, _ := claims.GetStr("email")
	firstName, _ := claims.GetStr("given_name")
	lastName, _ := claims.GetStr("family_name")
	picture, _ := claims.GetStr("picture")

	// Create user object to create new user or update existing user
	userData := models.User{
		Email:                 email,
		FirstName:             firstName,
		LastName:              lastName,
		Picture:               picture,
		AccessToken:           res.AccessToken,
		AccessTokenExpireDate: primitive.NewDateTimeFromTime(accessTokenExpireDate),
		RefreshToken:          res.RefreshToken,
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
			panic(err)
		}

		userId = res.InsertedID.(primitive.ObjectID)
	} else {
		// User does exist, get user id
		var user models.User
		if err := updateResult.Decode(&user); err != nil {
			panic(err)
		}

		userId = user.Id
	}

	// Set session variables
	session := sessions.Default(c)
	session.Set("userId", userId.Hex())
	session.Set("accessToken", res.AccessToken)
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
