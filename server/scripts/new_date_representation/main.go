package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"schej.it/server/db"
	"schej.it/server/models"
)

func main() {
	closeConnection := db.Init()
	defer closeConnection()

	cur, err := db.EventsCollection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.Background())

	type OldEvent struct {
		Id        primitive.ObjectID  `json:"_id" bson:"_id,omitempty"`
		OwnerId   primitive.ObjectID  `json:"ownerId" bson:"ownerId,omitempty"`
		Name      string              `json:"name" bson:"name,omitempty"`
		StartDate *primitive.DateTime `json:"startDate" bson:"startDate,omitempty"`
		EndDate   *primitive.DateTime `json:"endDate" bson:"endDate,omitempty"`

		// StartTime and EndTime are UTC hours, dates are an array of utc dates
		StartTime *float32 `json:"startTime" bson:"startTime,omitempty"`
		EndTime   *float32 `json:"endTime" bson:"endTime,omitempty"`
		Dates     []string `json:"dates" bson:"dates,omitempty"`

		Responses map[string]*models.Response `json:"responses" bson:"responses"`
	}

	var events []OldEvent
	for cur.Next(context.Background()) {
		var oldEvent OldEvent
		if err := cur.Decode(&oldEvent); err != nil {
			fmt.Println("couldn't decode event")
		}

		events = append(events, oldEvent)
	}

	for _, oldEvent := range events {
		var dates []primitive.DateTime
		var duration float32

		if oldEvent.StartDate != nil {
			startTime := oldEvent.StartDate.Time().UTC().Hour()
			endTime := oldEvent.EndDate.Time().UTC().Hour()
			duration = float32(endTime - startTime)
			if duration < 0 {
				duration += 24
			}

			curDate := oldEvent.StartDate.Time()
			for curDate.Before(oldEvent.EndDate.Time()) {
				dates = append(dates, primitive.NewDateTimeFromTime(curDate))
				curDate = curDate.Add(24 * time.Hour)
			}

			fmt.Println(duration, dates)
		} else if oldEvent.StartTime != nil {
			duration = *oldEvent.EndTime - *oldEvent.StartTime
			if duration < 0 {
				duration += 24
			}

			for _, dateString := range oldEvent.Dates {
				split := strings.Split(dateString, "-")
				year, _ := strconv.Atoi(split[0])
				month, _ := strconv.Atoi(split[1])
				day, _ := strconv.Atoi(split[2])

				date := time.Date(year, time.Month(month), day, int(*oldEvent.StartTime), 0, 0, 0, time.UTC)
				// fmt.Println(date)

				dates = append(dates, primitive.NewDateTimeFromTime(date))
			}

		}

		if dates != nil && duration > 0 {
			fmt.Println(oldEvent.Id)
			result, err := db.EventsCollection.UpdateOne(
				context.Background(),
				bson.M{"_id": oldEvent.Id},
				bson.M{
					"$set": bson.M{
						"dates":    dates,
						"duration": duration,
					},
				},
			)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(result.MatchedCount)
		}
	}
}
