package main

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"strings"
)

func cmdBalance(s *discordgo.Session, m *discordgo.Message, args []string) error {
	// User for balance
	var user *User
	var username string
	// Walrus operator won't work properly, so need to declare error ahead of time
	var err error
	if len(args) > 0 {
		member, err := searchGuild(s, m, strings.Join(args, " "))
		if err != nil {
			return err
		}

		if member.User.Bot {
			return errors.New("Specified user is a bot")
		}

		user, err = GetUser(member.User.ID)
		if err != nil {
			return err
		}

		username = member.User.Username
	} else {
		// Get current user's balance
		user, err = GetUser(m.Author.ID)
		if err != nil {
			return err
		}

		username = m.Author.Username
	}

	var stocksValue float64 = 0

	for symbol, count := range user.Shares {
		oneShareValue, err := getLivePrice(symbol, false)

		if err != nil {
			return err
		}

		stocksValue += oneShareValue * float64(count)
	}

	_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content: "",
		Embed: &discordgo.MessageEmbed{
			Title: username + "'s Account",
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:   "Net Worth",
					Value:  usdFormatter.FormatMoney(user.Balance + stocksValue),
					Inline: false,
				},
				&discordgo.MessageEmbedField{
					Name:   "Cash",
					Value:  usdFormatter.FormatMoney(user.Balance),
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Stocks",
					Value:  usdFormatter.FormatMoney(stocksValue),
					Inline: true,
				},
			},
			Color: 0x3E606F,
		},
	})

	return err
}
