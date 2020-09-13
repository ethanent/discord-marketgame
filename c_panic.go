package main

import (
	"errors"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func cmdPanic(s *discordgo.Session, m *discordgo.Message, args []string) error {
	if len(args) < 1 {
		return errors.New("Too few arguments.\nUsage: !panic <symbol | *>")
	}

	u, err := GetUser(m.Author.ID)
	if err != nil {
		return err
	}

	if args[0] == "*" || args[0] == "." {
		failures := map[string]error{}

		for symbol := range u.Shares {
			err = cmdPanic(s, m, []string{symbol})
			if err != nil {
				failures[symbol] = err
			}
		}

		if len(failures) > 0 {
			displayFailFields := []*discordgo.MessageEmbedField{}

			for symbol, err := range failures {
				displayFailFields = append(displayFailFields, &discordgo.MessageEmbedField{
					Name:  symbol,
					Value: ":x: " + err.Error(),
				})
			}

			s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
				Content: "",
				Embed: &discordgo.MessageEmbed{
					Title:  "Panic Partially Failed",
					Color:  0xFF0000,
					Fields: displayFailFields,
				},
			})
		}

		return nil
	} else {
		symbol := strings.ToUpper(args[0])

		shares, ok := u.Shares[symbol]
		if !ok {
			return errors.New("You don't own any " + symbol + ".")
		}

		return cmdSell(s, m, []string{symbol, strconv.Itoa(shares)})
	}
}
