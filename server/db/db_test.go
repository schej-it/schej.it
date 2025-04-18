package db

import (
	"fmt"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetDailyUserLogByDate(t *testing.T) {
	GetDailyUserLogByDate(time.Now(), 7)
}

func TestGenerateShortEventId(t *testing.T) {
	Init()

	objectId, _ := primitive.ObjectIDFromHex("6607d6409f96021811c0a55f")
	id := GenerateShortEventId(objectId)
	fmt.Println(id)
}
