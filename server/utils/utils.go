package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/brianvoe/sjwt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"schej.it/server/logger"
	"schej.it/server/models"
)

// Returns whether running on production server
func IsRelease() bool {
	mode := os.Getenv("GIN_MODE")
	return mode == "release"
}

func PrintJson(s interface{}) {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		logger.StdErr.Panicln(err)
	}

	fmt.Println(string(data))
}

func ParseJWT(jwt string) sjwt.Claims {
	claims, err := sjwt.Parse(jwt)
	if err != nil {
		logger.StdErr.Panicln(err)
	}

	return claims
}

func StringToObjectID(s string) primitive.ObjectID {
	objectID, err := primitive.ObjectIDFromHex(s)
	if err != nil {
		logger.StdErr.Panicln(err)
	}

	return objectID
}

// Returns the currently signed in user
func GetAuthUser(c *gin.Context) *models.User {
	userInterface, _ := c.Get("authUser")
	user := userInterface.(*models.User)
	return user
}

// Gets the access token expire date from an "expiresIn" int representing the number of seconds
// after which the access token will expire
func GetAccessTokenExpireDate(expiresIn int) time.Time {
	expireDuration, err := time.ParseDuration(fmt.Sprintf("%ds", expiresIn))
	if err != nil {
		logger.StdErr.Panicln(err)
	}
	return time.Now().Add(expireDuration)
}

// Returns the ISO date string for the given date
func GetDateString(date time.Time) string {
	s, _ := date.UTC().MarshalText()
	return string(s)[:10]
}

// Returns a time object with the given date and a time string in the form of "00:00:00"
func GetDateAtTime(date time.Time, timeString string) time.Time {
	utcDateString := GetDateString(date)
	newDate, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT%sZ", utcDateString, timeString))
	if err != nil {
		logger.StdErr.Panicln(err)
	}
	return newDate
}

// Escapes regex for a string
func EscapeRegExp(str string) string {
	check := regexp.MustCompile(`([.*+?^${}()|[\]\\])`)
	return check.ReplaceAllString(str, "\\${1}")
}

// Returns the correct client id given the token origin
func GetClientIdFromTokenOrigin(tokenOrigin models.TokenOriginType) string {
	switch tokenOrigin {
	case models.ANDROID:
		return os.Getenv("ANDROID_CLIENT_ID")
	case models.IOS:
		return os.Getenv("IOS_CLIENT_ID")
	default:
		return os.Getenv("CLIENT_ID")
	}
}

// Prints the http response as a string
func PrintHttpResponse(resp *http.Response) {
	body, _ := io.ReadAll(resp.Body)
	logger.StdOut.Println(string(body))
	resp.Body = io.NopCloser(bytes.NewBuffer(body))
}
