package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"schej.it/server/db"
	"schej.it/server/logger"
	"schej.it/server/models"
	"schej.it/server/utils"
)

var activeUsers Command = Command{
	Name: "!active_users",
	Description: `Gets the number of active users in the database, based on last sign in date. 
  - if LIST is true, it will list the name/email of all users, otherwise, it will show a bar graph
  - DAYS is the amount of days since last sign in
  `,
	Usage: "!active_users [LIST=false] [DAYS=7]",
	Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
		var err error

		// Parse args
		list := false
		days := 7
		if len(args) >= 1 {
			if args[0] == "true" {
				list = true
			} else if args[0] == "false" {
				list = false
			} else {
				sendMessage(s, m, fmt.Sprintf("LIST=%s is not a valid boolean!", args[0]))
				return
			}
		}
		if len(args) >= 2 {
			days, err = strconv.Atoi(args[1])
			if err != nil {
				sendMessage(s, m, fmt.Sprintf("DAYS=%s is not a valid number!", args[1]))
				return
			}
		}

		// Query for daily user logs starting from `days` days before the current date
		startDate := time.Now().AddDate(0, 0, -days)
		startDate = utils.GetDateAtTime(startDate, "00:00:00")
		query := bson.M{"date": bson.M{"$gte": primitive.NewDateTimeFromTime(startDate)}}
		sort := bson.M{"date": -1}

		var logs []models.DailyUserLog
		if list {
			// Find and populate
			cursor, err := db.DailyUserLogCollection.Aggregate(context.Background(), []bson.M{
				{"$match": query},
				{"$sort": sort},
				{"$lookup": bson.M{
					"from":         "users",
					"localField":   "userIds",
					"foreignField": "_id",
					"as":           "users",
				}},
				{"$project": bson.M{
					"date":            1,
					"users._id":       1,
					"users.firstName": 1,
					"users.lastName":  1,
					"users.email":     1,
				}},
			})
			if err != nil {
				logger.StdErr.Panicln(err)
			}
			if err := cursor.All(context.Background(), &logs); err != nil {
				logger.StdErr.Panicln(err)
			}
		} else {
			// Find matches
			cursor, err := db.DailyUserLogCollection.Find(context.Background(), query, &options.FindOptions{
				Sort: sort,
			})
			if err != nil {
				logger.StdErr.Panicln(err)
			}
			if err := cursor.All(context.Background(), &logs); err != nil {
				logger.StdErr.Panicln(err)
			}
		}

		// Add empty days
		curDate := startDate
		for i := len(logs) - 1; i >= 0; i-- {
			// Add all dates up to the current log date
			for !logs[i].Date.Time().Equal(curDate) && curDate.Before(time.Now()) {
				// Insert curDate into logs, with an empty users array
				logs, err = utils.Insert(logs, i+1, models.DailyUserLog{
					Date:  primitive.NewDateTimeFromTime(curDate),
					Users: make([]models.User, 0),
				})
				if err != nil {
					logger.StdErr.Panicln(err)
				}
				curDate = curDate.AddDate(0, 0, 1)
			}

			// Increase curDate by a day
			curDate = curDate.AddDate(0, 0, 1)
		}

		// Add all dates up to the current date
		for curDate.Before(time.Now()) {
			logs, err = utils.Insert(logs, 0, models.DailyUserLog{
				Date:  primitive.NewDateTimeFromTime(curDate),
				Users: make([]models.User, 0),
			})
			if err != nil {
				logger.StdErr.Panicln(err)
			}
			curDate = curDate.AddDate(0, 0, 1)
		}

		// Define constants
		dayStrings := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}

		if list {
			// Display a list of all active users
			sendMessage(s, m, "Active Users:\n")
			message := ""
			for _, log := range logs {
				date := log.Date.Time()
				message += dayStrings[date.Weekday()] + " "
				message += utils.GetDateString(date) + " | "
				message += fmt.Sprintf("Count: %d\n", len(log.Users))

				for _, user := range log.Users {
					message += fmt.Sprintf("\t- %s %s (%s)\n", user.FirstName, user.LastName, user.Email)
				}
			}

			for _, msg := range splitLongMessage(message, "```") {
				sendMessage(s, m, msg)
			}
		} else {
			// Display a bar graph of active users over time

			// Generate labels and data based on logs
			labels := make([]string, 0)
			data := make([]int, 0)
			for i := len(logs) - 1; i >= 0; i-- {
				labels = append(labels, utils.GetDateString(logs[i].Date.Time()))
				data = append(data, len(logs[i].UserIds))
			}

			// Generate chart using QuickChart API
			chart := bson.M{
				"type": "bar",
				"data": bson.M{
					"labels": labels,
					"datasets": bson.A{bson.M{
						"label": "Active Users",
						"data":  data,
					}},
				},
				"options": bson.M{
					"scales": bson.M{
						"yAxes": bson.A{bson.M{
							"ticks": bson.M{
								"stepSize": 1,
							},
						}},
					},
				},
			}
			jsonStr, _ := json.Marshal(chart)

			encodedChart := url.PathEscape(string(jsonStr))
			chartUrl := fmt.Sprintf(`https://quickchart.io/chart?c=%s&backgroundColor=white`, encodedChart)
			chartEmbed := &discordgo.MessageEmbed{
				Title: "Active Users",
				Image: &discordgo.MessageEmbedImage{
					URL: chartUrl,
				},
			}

			sendEmbed(s, m, chartEmbed)
		}
	},
}
