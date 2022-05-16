package utils

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/brianvoe/sjwt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
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

// Returns a document if it exists, otherwise returns nil
func FindDocumentIfExists(coll *mongo.Collection, filter interface{}) *mongo.SingleResult {
	count, err := coll.CountDocuments(context.Background(), filter)
	if err != nil {
		panic(err)
	}

	if count == 0 {
		return nil
	}

	return coll.FindOne(context.Background(), filter)
}
