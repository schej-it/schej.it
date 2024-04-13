package slackbot

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"schej.it/server/models"
	"schej.it/server/slackbot/commands"
	"schej.it/server/utils"
)

func SendEventCreatedMessage(insertedId string, creator string, event models.Event) {
	eventInfoText := fmt.Sprintf(
		"*Name*: %s\n"+
			"*Event url*: https://schej.it/e/%s\n"+
			"*Short url*: https://schej.it/e/%s\n"+
			"*Creator*: %s\n"+
			"*Num days*: %v\n"+
			"*Type*: %s",
		event.Name,
		insertedId,
		utils.Coalesce(event.ShortId),
		creator,
		len(event.Dates),
		event.Type,
	)

	if event.Type == models.GROUP {
		eventInfoText += fmt.Sprintf("\n*Num attendees*: %v", len(utils.Coalesce(event.Attendees)))
	} else {
		eventInfoText += fmt.Sprintf(
			"\n*Notifications Enabled*: %v\n"+
				"*Num remindees*: %v",
			utils.Coalesce(event.NotificationsEnabled),
			len(utils.Coalesce(event.Remindees)),
		)
	}

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
				"text": eventInfoText,
			},
		},
	}}

	SendMessage(&response)
}
