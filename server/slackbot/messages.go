package slackbot

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"schej.it/server/models"
	"schej.it/server/slackbot/commands"
)

func SendEventCreatedMessage(insertedId string, creator string, event models.Event) {
	response := commands.Response{Blocks: []bson.M{
		{
			"type": "header",
			"text": bson.M{
				"type":  "plain_text",
				"text":  ":tada: New event created! :tada:",
				"emoji": true,
			},
		},
		{
			"type": "section",
			"text": bson.M{
				"type": "mrkdwn",
				"text": fmt.Sprintf(
					"*Event url*: https://schej.it/e/%s\n"+
						"*Creator*: %s\n"+
						"*Type*: %s\n"+
						"*Notifications Enabled*: %v\n"+
						"*Num remindees*: %v",
					insertedId,
					creator,
					event.Type,
					event.NotificationsEnabled,
					len(event.Remindees),
				),
			},
		},
	}}

	SendMessage(&response)
}
