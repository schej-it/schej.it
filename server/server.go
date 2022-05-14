package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	router := gin.Default()

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

	// Auth routes
	auth := router.Group("/auth")
	{
		auth.POST("/sign-in", signIn)
	}

	// Run server
	router.Run(":3000")
}

func loadDotEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

type signInPayload struct {
	Code string `json:"code" binding:"required"`
}

func signIn(c *gin.Context) {
	var payload signInPayload

	if err := c.BindJSON(&payload); err != nil {
		return
	}

	/*formattedUrl := fmt.Sprintf(
		"?client_id=%s&client_secret=%s&code=%s&grant_type=authorization_code&redirect_uri=%s",
		os.Getenv("CLIENT_ID"),
		os.Getenv("CLIENT_SECRET"),
		payload.Code,
		url.QueryEscape("http://localhost:8080"),
	)

	resp, err := http.Post(formattedUrl, )*/

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

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))

	c.JSON(http.StatusOK, gin.H{"test": payload.Code})
}
