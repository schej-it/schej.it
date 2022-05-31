package discord_bot

import (
	"os"
	"regexp"

	"github.com/bwmarrin/discordgo"
	"schej.it/server/discord_bot/commands"
	"schej.it/server/logger"
	"schej.it/server/utils"
)

var BotId string
var bot *discordgo.Session
var listeningChannel *discordgo.Channel
var commandMap map[string]commands.Command = make(map[string]commands.Command)

// Initialize the discord bot
func Init() {
	var err error
	token := os.Getenv("DISCORD_BOT_TOKEN")
	guildId := os.Getenv("GUILD_ID")
	bot, err = discordgo.New("Bot " + token)
	if err != nil {
		logger.StdErr.Panicln(err)
	}

	u, err := bot.User("@me")
	if err != nil {
		logger.StdErr.Panicln(err)
	}
	BotId = u.ID

	err = bot.Open()
	if err != nil {
		logger.StdErr.Panicln(err)
	}
	logger.StdOut.Println("Discord bot initialized")

	// Get the channel object to listen on
	var listeningChannelName string
	if utils.IsRelease() {
		listeningChannelName = "schej-it-bot"
	} else {
		listeningChannelName = "schej-it-bot-dev"
	}
	channels, _ := bot.GuildChannels(guildId)
	for _, channel := range channels {
		if channel.Name == listeningChannelName {
			listeningChannel = channel
			break
		}
	}

	// Construct commandMap
	commands.Init()
	for _, command := range commands.Commands {
		commandMap[command.Name] = command
	}

	// Bot handlers
	bot.AddHandler(messageHandler)
}

// Send a message to the listening channel
func SendMessage(message string) {
	_, _ = bot.ChannelMessageSend(listeningChannel.ID, message)
}

// messageHandler is run every time a message is sent in a server we are listening to
func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.ChannelID != listeningChannel.ID {
		return
	}
	if m.Author.ID == BotId {
		return
	}

	// Parse the command + arguments
	re := regexp.MustCompile(" +")
	args := re.Split(m.Content, -1)
	commandName := args[0]
	args = args[1:]

	// Execute command if it exists
	if command, ok := commandMap[commandName]; ok {
		command.Execute(s, m, args)
	}
}
