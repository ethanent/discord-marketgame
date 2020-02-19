package main

import (
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func cmdCancel(s *discordgo.Session, m *discordgo.Message, args []string) error {
	if len(args) < 1 {
		return errors.New("Missing arguments.\nUsage: !cancel <symbol>")
	}

	symbol := strings.ToUpper(args[0])

	// Get user

	u, err := GetUser(m.Author.ID)

	if err != nil {
		return err
	}

	// Check that stop order exists

	_, ok := u.Stops[symbol]

	if !ok {
		return errors.New("The stop order does not exist.")
	}

	// Cancel the order

	delete(u.Stops, symbol)

	// Save

	err = u.Save()

	if err != nil {
		return err
	}

	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content: "",
		Embed: &discordgo.MessageEmbed{
			Title: ":ballot_box_with_check: " + symbol + " Stop Order Cancelled",
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:  "Information",
					Value: "Your stop order has been cancelled and the lock on the symbol has been lifted.",
				},
			},
			Color: 0x46E8B2,
		},
	})

	return nil
}
