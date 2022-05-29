package commands

import "github.com/bwmarrin/discordgo"

var activeUsers Command = Command{
	Name: "!active_users",
	Description: `Gets the number of active users in the database, based on last sign in date. 
  - if LIST is true, it will list the name/email of all users, otherwise, it will show a bar graph
  - DAYS is the amount of days since last sign in
  `,
	Usage: "!active_users [LIST=false] [DAYS=7]",
	Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
		sendMessage(s, m, "active users was called hehe ")
	},
}
