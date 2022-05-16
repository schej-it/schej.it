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
		authRouter.GET("/status", middleware.AuthRequired(), getAuthStatus)
	}

	// Current user routes
	userRouter := router.Group("/user")
	userRouter.Use(middleware.AuthRequired())
	{
		userRouter.GET("/availability", getAvailability)
		userRouter.GET("/events", getUserEvents)
		userRouter.GET("/profile", getUserProfile)
	}

	// Event routes
	eventRouter := router.Group("/events")
	{
		eventRouter.POST("", middleware.AuthRequired(), createEvent)
		eventRouter.GET("/:eventId", middleware.AuthRequired(), getEvent)
		eventRouter.POST("/:eventId/response", middleware.AuthRequired(), updateResponse)
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

		userId = user.Id
	}

	// Set session variables
	session := sessions.Default(c)
	session.Set("userId", userId.Hex())
	session.Set("accessToken", res["access_token"])
	session.Save()

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// Gets whether the user is signed in or not
func getAuthStatus(c *gin.Context) {
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

// Gets all the user's events
func getUserEvents(c *gin.Context) {
	session := sessions.Default(c)

	var events []models.Event
	cursor, err := db.EventsCollection.Find(context.Background(), bson.M{
		"ownerId": utils.GetUserId(session),
	})
	if err != nil {
		panic(err)
	}
	if err := cursor.All(context.Background(), &events); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, events)
}

// Gets the user's profile
func getUserProfile(c *gin.Context) {
	user, _ := c.Get("authUser")
	user = user.(*models.User)

	c.JSON(http.StatusOK, user)
}

// Creates a new event
func createEvent(c *gin.Context) {
	payload := struct {
		Name      string    `json:"name" binding:"required"`
		StartDate time.Time `json:"startDate" binding:"required"`
		EndDate   time.Time `json:"endDate" binding:"required"`
		StartTime int       `json:"startTime" binding:"required"`
		EndTime   int       `json:"endTime" binding:"required"`
	}{}
	if err := c.Bind(&payload); err != nil {
		return
	}
	session := sessions.Default(c)

	event := models.Event{
		OwnerId:   utils.GetUserId(session),
		Name:      payload.Name,
		StartDate: primitive.NewDateTimeFromTime(payload.StartDate),
		EndDate:   primitive.NewDateTimeFromTime(payload.EndDate),
		StartTime: payload.StartTime,
		EndTime:   payload.EndTime,
	}

	result, err := db.EventsCollection.InsertOne(context.Background(), event)
	if err != nil {
		panic(err)
	}

	insertedId := result.InsertedID.(primitive.ObjectID).Hex()
	c.JSON(http.StatusCreated, gin.H{"eventId": insertedId})
}

// Gets an event located at the given id
func getEvent(c *gin.Context) {
	eventId := c.Param("eventId")
	event := db.GetEventById(eventId)

	c.JSON(http.StatusOK, event)
}

// Updates the current user's availability
func updateResponse(c *gin.Context) {
	payload := struct {
		Availability []string `json:"availability" binding:"required"`
	}{}
	if err := c.Bind(&payload); err != nil {
		return
	}
	session := sessions.Default(c)

	eventId := c.Param("eventId")

	response := models.Response{
		UserId:       utils.GetUserId(session),
		Availability: payload.Availability,
	}

	userIdString := session.Get("userId").(string)
	_, err := db.EventsCollection.UpdateByID(
		context.Background(),
		utils.StringToObjectID(eventId),
		bson.M{
			"$set": bson.M{
				"responses." + userIdString: response,
			},
		},
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("YAY!")

	c.JSON(http.StatusOK, gin.H{"success": true})
}
