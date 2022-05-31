package commands

import "github.com/bwmarrin/discordgo"

type Command struct {
	Name        string
	Description string
	Usage       string
	Execute     func(s *discordgo.Session, m *discordgo.MessageCreate, args []string)
}

var Commands []Command = make([]Command, 0)

func Init() {
	Commands = append(Commands, activeUsers)
	Commands = append(Commands, help)
}

// Send a message to the current channel that received a message
func sendMessage(s *discordgo.Session, m *discordgo.MessageCreate, message string) {
	_, _ = s.ChannelMessageSend(m.ChannelID, message)
}
