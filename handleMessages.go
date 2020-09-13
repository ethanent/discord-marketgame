package main

import (
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// CommandHandler is a function used to handle commands
type CommandHandler func(*discordgo.Session, *discordgo.Message, []string) error

var commandHandlers map[string]CommandHandler = map[string]CommandHandler{}

var alternateCommands map[string]string = map[string]string{}

func registerCommand(command string, h CommandHandler) {
	commandHandlers[command] = h
}

func registerAlternate(from string, to string) {
	alternateCommands[from] = to
}

func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg := *(*m).Message

	contentChars := []rune(msg.Content)

	if len(contentChars) > 0 && string(contentChars[0]) == "!" {
		args := strings.Split(string([]rune(msg.Content)[1:]), " ")

		lookupCommand := args[0]

		alternateFor, ok := alternateCommands[args[0]]

		if ok {
			lookupCommand = alternateFor
		}

		handler, ok := commandHandlers[lookupCommand]

		if ok {
			err := handler(s, &msg, args[1:])

			if err != nil {
				os.Stderr.WriteString(err.Error())

				_, err = s.ChannelMessageSendComplex(msg.ChannelID, &discordgo.MessageSend{
					Content: "",
					Embed: &discordgo.MessageEmbed{
						Title:       ":x: Error",
						Description: err.Error(),
						Color:       0xFF0000,
					},
				})
			}
		}
	}
}

func displayError(s *discordgo.Session, m *discordgo.Message, err error) {
	os.Stderr.WriteString(err.Error())

	_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content: "",
		Embed: &discordgo.MessageEmbed{
			Title:       "ERROR",
			Description: err.Error(),
			Color:       0xFF0000,
		},
	})
}