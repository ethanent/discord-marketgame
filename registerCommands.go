package main

import (
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	av "github.com/cmckee-dev/go-alpha-vantage"
)

func registerCommands() {
	registerCommand("help", func(s *discordgo.Session, m *discordgo.Message, args []string) error {
		_, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content: "",
			Embed: &discordgo.MessageEmbed{
				Title:       "Help",
				Description: "Welcome to MarketBot.\n\nCommands:\n\n!help",
				Color:       0x3E606F,
			},
		})

		return err
	})

	registerCommand("bal", func(s *discordgo.Session, m *discordgo.Message, args []string) error {
		user, err := getUser(m.Author.ID)

		if err != nil {
			_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
				Content: "",
				Embed: &discordgo.MessageEmbed{
					Title:       "ERROR",
					Description: err.Error(),
					Color:       0xFF0000,
				},
			})

			return err
		}

		_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content: "",
			Embed: &discordgo.MessageEmbed{
				Title: m.Author.Username + "'s Account",
				Fields: []*discordgo.MessageEmbedField{
					&discordgo.MessageEmbedField{
						Name:  "Balance",
						Value: usdFormatter.FormatMoney(user.Balance),
					},
				},
				Color: 0x3E606F,
			},
		})

		return err
	})

	registerCommand("price", func(s *discordgo.Session, m *discordgo.Message, args []string) error {
		symbol := strings.ToUpper(args[1])

		tsvs, err := avClient.StockTimeSeries(av.TimeSeriesDaily, symbol)

		if err != nil {
			_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
				Content: "",
				Embed: &discordgo.MessageEmbed{
					Title:       "ERROR",
					Description: err.Error(),
					Color:       0xFF0000,
				},
			})

			return err
		}

		price := *tsvs[len(tsvs)-1]
		yesterdayPrice := *tsvs[len(tsvs)-2]

		pxChange := (price.Close/yesterdayPrice.Close)*100 - 100
		addSymbol := ""

		if pxChange > 0 {
			addSymbol = "+"
		}

		_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content: "",
			Embed: &discordgo.MessageEmbed{
				Title: "$" + symbol + " Stock",
				Fields: []*discordgo.MessageEmbedField{
					&discordgo.MessageEmbedField{
						Name:  "Share price",
						Value: usdFormatter.FormatMoney(price.Close) + " (" + addSymbol + strconv.FormatFloat(pxChange, 'f', 2, 64) + "% from yesterday)",
					},
				},
				Color: 0x3E606F,
			},
		})

		return err
	})
}
