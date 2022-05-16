package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"encoding/json"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"schej.it/server/db"
	"schej.it/server/middleware"
	"schej.it/server/models"
	"schej.it/server/utils"
)

func main() {
	router := gin.Default()

	// Load .env variables
	loadDotEnv()

	// Cors
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PATCH"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Init database
	closeConnection := db.Init()
	defer closeConnection()

	// Session
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("session", store))

	// Auth routes
	authRouter := router.Group("/auth")
	{
		authRouter.POST("/sign-in", signIn)
	}

	// Current user routes
	userRouter := router.Group("/user")
	userRouter.Use(middleware.AuthRequired())
	{
		userRouter.GET("/availability", getAvailability)
	}

	// Event routes
	eventRouter := router.Group("/event/")
	{
		eventRouter.POST("/", middleware.AuthRequired(), createEvent)
	}

	// Run server
	router.Run(":3000")
}

// Load .env variables
func loadDotEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		panic("Error loading .env file")
	}
}

// Sign user in and set the access token session variable
func signIn(c *gin.Context) {
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

		userId = user.ID
	}

	// Set session variables
	session := sessions.Default(c)
	session.Set("userId", userId.Hex())
	session.Set("accessToken", res["access_token"])
	session.Save()

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// gets the times that the current user is available
func getAvailability(c *gin.Context) {
	// Bind query parameters
	payload := struct {
		TimeMin string `form:"timeMin" binding:"required"`
		TimeMax string `form:"timeMax" binding:"required"`
	}{}
	if err := c.Bind(&payload); err != nil {
		return
	}
	session := sessions.Default(c)

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("https://www.googleapis.com/calendar/v3/calendars/primary/events?timeMin=%s&timeMax=%s", payload.TimeMin, payload.TimeMax),
		nil,
	)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", session.Get("accessToken")))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)

	c.JSON(http.StatusOK, res)
}

// Creates a new event
func createEvent(c *gin.Context) {
	payload := struct {
		Name      string `json:"name" binding:"required"`
		StartDate string `json:"startDate" binding:"required"`
		EndDate   string `json:"endDate" binding:"required"`
		StartTime string `json:"startTime" binding:"required"`
		EndTime   string `json:"endTime" binding:"required"`
	}{}
	if err := c.Bind(&payload); err != nil {
		return
	}

}
