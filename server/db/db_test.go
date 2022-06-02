package db

import (
	"testing"
	"time"
)

func TestGetDailyUserLogByDate(t *testing.T) {
	GetDailyUserLogByDate(time.Now(), 7)
}
