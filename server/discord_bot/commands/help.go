package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var help Command = Command{
	Name:        "!help",
	Description: "Displays this help message",
	Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
		var message string
		for _, command := range Commands {
			message += fmt.Sprintf("\n\n`%s`: %s ", command.Name, command.Description)
			if len(command.Usage) > 0 {
				message += fmt.Sprintf("Usage: `%s`", command.Usage)
			}
		}

		sendMessage(s, m, message)
	},
}
