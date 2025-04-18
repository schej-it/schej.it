package db_test

import (
	"fmt"
	"testing"
	"time"

	"schej.it/server/db"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetDailyUserLogByDate(t *testing.T) {
	db.GetDailyUserLogByDate(time.Now(), 7)
}

func TestGenerateShortEventId(t *testing.T) {
	db.Init()

	objectId, _ := primitive.ObjectIDFromHex("6607d6409f96021811c0a55f")
	id := db.GenerateShortEventId(objectId)
	fmt.Println(id)
}
