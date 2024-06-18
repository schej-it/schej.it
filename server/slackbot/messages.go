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
			"*Type*: %s\n",
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
		daysOnly := utils.Coalesce(event.DaysOnly)
		notificationsEnabled := utils.Coalesce(event.NotificationsEnabled)
		numRemindees := len(utils.Coalesce(event.Remindees))
		blindAvailabilityEnabled := utils.Coalesce(event.BlindAvailabilityEnabled)
		sendEmailAfterXResponses := utils.Coalesce(event.SendEmailAfterXResponses)
		when2meetHref := utils.Coalesce(event.When2meetHref)

		eventInfoText += fmt.Sprintln("*Days only*:", daysOnly)
		if notificationsEnabled {
			eventInfoText += fmt.Sprintln("*Notifications enabled*:", notificationsEnabled)
		}
		if numRemindees > 0 {
			eventInfoText += fmt.Sprintln("*Num remindees*:", numRemindees)
		}
		if blindAvailabilityEnabled {
			eventInfoText += fmt.Sprintln("*Blind availability*:", blindAvailabilityEnabled)
		}
		if sendEmailAfterXResponses > 0 {
			eventInfoText += fmt.Sprintln("*Send email after X responses*:", sendEmailAfterXResponses)
		}
		if len(when2meetHref) > 0 {
			eventInfoText += fmt.Sprintf("*When2meet URL*: https://when2meet.com%s\n", when2meetHref)
		}
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
