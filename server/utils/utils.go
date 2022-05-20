package utils

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/brianvoe/sjwt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func PrintJson(s gin.H) {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}

func ParseJWT(jwt string) sjwt.Claims {
	claims, err := sjwt.Parse(jwt)
	if err != nil {
		panic(err)
	}

	return claims
}

func StringToObjectID(s string) primitive.ObjectID {
	objectID, err := primitive.ObjectIDFromHex(s)
	if err != nil {
		panic(err)
	}

	return objectID
}

// Gets the user id from the current session as an ObjectID object
func GetUserId(session sessions.Session) primitive.ObjectID {
	return StringToObjectID(session.Get("userId").(string))
}

// Gets the access token expire date from an "expiresIn" int representing the number of seconds
// after which the access token will expire
func GetAccessTokenExpireDate(expiresIn int) time.Time {
	expireDuration, err := time.ParseDuration(fmt.Sprintf("%ds", expiresIn))
	if err != nil {
		panic(err)
	}
	return time.Now().Add(expireDuration)
}
