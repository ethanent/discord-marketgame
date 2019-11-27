package main

import (
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// CommandHandler is a function used to handle commands
type CommandHandler func(*discordgo.Session, *discordgo.Message, []string) error

var commandHandlers map[string]CommandHandler = map[string]CommandHandler{}

func registerCommand(command string, h CommandHandler) {
	commandHandlers[command] = h
}

func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg := *(*m).Message

	contentChars := []rune(msg.Content)

	if len(contentChars) > 0 && string(contentChars[0]) == "!" {
		args := strings.Split(string([]rune(msg.Content)[1:]), " ")

		handler, ok := commandHandlers[args[0]]

		if ok {
			err := handler(s, &msg, args[1:])

			if err != nil {
				os.Stderr.WriteString(err.Error())

				_, err = s.ChannelMessageSendComplex(msg.ChannelID, &discordgo.MessageSend{
					Content: "",
					Embed: &discordgo.MessageEmbed{
						Title:       "ERROR",
						Description: err.Error(),
						Color:       0xFF0000,
					},
				})
			}
		}
	}
}
