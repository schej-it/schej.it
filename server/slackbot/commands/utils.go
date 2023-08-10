package commands

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"schej.it/server/logger"
)

// Blocks documentation: https://api.slack.com/reference/block-kit/blocks
// Block builder: https://app.slack.com/block-kit-builder/

type Response struct {
	ResponseType string   `json:"response_type,omitempty"` // ephemeral or in_channel
	Text         string   `json:"text,omitempty"`
	Blocks       []bson.M `json:"blocks,omitempty"`
}

type Command struct {
	Name        string
	Description string
	Usage       string
	Execute     func(args []string, webhookUrl string)
}

var CommandMap = map[string]Command{
	activeUsers.Name: activeUsers,
	numUsers.Name:    numUsers,
}

func SendRawMessage(message *Response, webhookUrl string) {
	bodyBytes, _ := json.Marshal(message)
	bodyBuffer := bytes.NewBuffer(bodyBytes)

	req, err := http.NewRequest("POST", webhookUrl, bodyBuffer)
	if err != nil {
		logger.StdErr.Println("Failed to send message to slack bot: ", err)
		return
	}
	req.Header.Set("Content-type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.StdErr.Println("Failed to send message to slack bot: ", err)
		return
	}
	defer resp.Body.Close()
}

// Returns an array containing the separate messages to send when a
// single message is too long (>2000 characters).
// Each message in the array will be surrounded by surroundString (e.g. "```")
func splitLongMessage(message string, surroundString string) []string {

	charLimit := 2000
	charLimit -= len(surroundString) * 2
	messageArray := make([]string, 0)

	// Go through message and keep substringing it in 2000 character chunks,
	// and add the chunks to the messageArray
	for len(message) > charLimit {
		// Split at the last newline before the character limit, or the character limit
		// if newline was not found
		splitIndex := strings.LastIndex(message[:charLimit], "\n")
		if splitIndex == -1 {
			splitIndex = charLimit
		}

		messageArray = append(messageArray, surroundString+message[:splitIndex]+surroundString)
		if message[splitIndex:splitIndex+1] == "\n" {
			splitIndex++
		}
		message = message[splitIndex:]
	}
	if len(message) > 0 {
		messageArray = append(messageArray, surroundString+message+surroundString)
	}

	return messageArray
}
