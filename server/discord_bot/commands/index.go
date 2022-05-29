package commands

import "github.com/bwmarrin/discordgo"

type Command struct {
	Name        string
	Description string
	Usage       string
	Execute     func(s *discordgo.Session, m *discordgo.MessageCreate, args []string)
}

var Commands []Command = []Command{
	activeUsers,
}

// Send a message to the current channel that received a message
func sendMessage(s *discordgo.Session, m *discordgo.MessageCreate, message string) {
	_, _ = s.ChannelMessageSend(m.ChannelID, message)
}
