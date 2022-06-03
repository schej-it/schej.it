package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Name        string
	Description string
	Usage       string
	Execute     func(s *discordgo.Session, m *discordgo.MessageCreate, args []string)
}

var Commands []Command = make([]Command, 0)

func Init() {
	Commands = append(Commands, activeUsers)
	Commands = append(Commands, numUsers)

	Commands = append(Commands, help)
}

// Send a message to the current channel that received a message
func sendMessage(s *discordgo.Session, m *discordgo.MessageCreate, message string) {
	_, _ = s.ChannelMessageSend(m.ChannelID, message)
}

func sendEmbed(s *discordgo.Session, m *discordgo.MessageCreate, embed *discordgo.MessageEmbed) {
	_, _ = s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func splitLongMessage(message string, surroundString string) []string {
	/* Returns an array containing the separate messages to send when a
	 * single message is too long (>2000 characters).
	 * Each message in the array will be surrounded by surroundString (e.g. "```")
	 */

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
