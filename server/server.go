package main

import (
	"fmt"
	"log"
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
	{
		userRouter.GET("/availability", getAvailability)
	}

	// Run server
	router.Run(":3000")
}

// Load .env variables
func loadDotEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
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

	session := sessions.Default(c)

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
		log.Fatal(err)
	}

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)

	fmt.Println(res["access_token"])
	session.Set("access_token", res["access_token"])
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

	fmt.Println("ACCESS TOKEN: ", session.Get("access_token"))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", session.Get("access_token")))
	fmt.Println(req)
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)

	c.JSON(http.StatusOK, res)
}
