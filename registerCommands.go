package main

import (
	"errors"
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
				Description: "Welcome to MarketGame. May the odds be in your favor.\n[Source](https://github.com/ethanent/discord-marketgame)",
				Footer: &discordgo.MessageEmbedFooter{
					Text: "(c) 2019 Ethan Davis",
				},
				Fields: []*discordgo.MessageEmbedField{
					&discordgo.MessageEmbedField{
						Name:   "Information Commands",
						Value:  "\n!help\n!price <ticker>\n!bal\n!shares",
						Inline: true,
					},
					&discordgo.MessageEmbedField{
						Name:   "Activity Commands",
						Value:  "\n!buy <ticker> <count>\n!reset",
						Inline: true,
					},
				},
				Color: 0x3E606F,
			},
		})

		return err
	})

	registerCommand("bal", func(s *discordgo.Session, m *discordgo.Message, args []string) error {
		user, err := GetUser(m.Author.ID)

		if err != nil {
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
		if len(args) < 1 {
			return errors.New("Missing arguments.\nUsage: !price <symbol>")
		}

		symbol := strings.ToUpper(args[0])

		tsvs, err := avClient.StockTimeSeries(av.TimeSeriesDaily, symbol)

		if err != nil {
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

	registerCommand("buy", func(s *discordgo.Session, m *discordgo.Message, args []string) error {
		if len(args) < 2 {
			return errors.New("Missing arguments.\nUsage: !buy <count> <symbol>")
		}

		count, err := strconv.Atoi(args[0])
		symbol := strings.ToUpper(args[1])

		if err != nil {
			return err
		}

		if count < 0 {
			return errors.New("You can't buy negative shares.")
		}

		sharePx, err := getLivePrice(symbol)

		if err != nil {
			return err
		}

		u, err := GetUser(m.Author.ID)

		if err != nil {
			return err
		}

		totalPx := float64(count) * sharePx

		if u.Balance-totalPx < 0 {
			return errors.New("You can't afford to buy " + strconv.Itoa(count) + "x$" + symbol)
		}

		u.Balance -= totalPx

		_, ok := u.Shares[symbol]

		if ok {
			u.Shares[symbol] += uint(count)
		} else {
			u.Shares[symbol] = uint(count)
		}

		err = u.Save()

		if err != nil {
			return err
		}

		_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content: "",
			Embed: &discordgo.MessageEmbed{
				Title: strconv.Itoa(count) + "x$" + symbol + " Purchase Complete",
				Fields: []*discordgo.MessageEmbedField{
					&discordgo.MessageEmbedField{
						Name:  "Buy price",
						Value: usdFormatter.FormatMoney(sharePx),
					},
					&discordgo.MessageEmbedField{
						Name:  "Total cost",
						Value: usdFormatter.FormatMoney(totalPx),
					},
				},
				Color: 0x3E606F,
			},
		})

		if err != nil {
			return err
		}

		return nil
	})

	registerCommand("sell", func(s *discordgo.Session, m *discordgo.Message, args []string) error {
		if len(args) < 2 {
			return errors.New("Missing arguments.\nUsage: !sell <count> <symbol>")
		}

		count, err := strconv.Atoi(args[0])
		symbol := strings.ToUpper(args[1])

		if err != nil {
			return err
		}

		if count < 0 {
			return errors.New("You can't sell negative shares.")
		}

		sharePx, err := getLivePrice(symbol)

		if err != nil {
			return err
		}

		u, err := GetUser(m.Author.ID)

		if err != nil {
			return err
		}

		totalPx := float64(count) * sharePx

		u.Balance += totalPx

		_, ok := u.Shares[symbol]

		if ok {
			if u.Shares[symbol] < uint(count) {
				return errors.New("You do not own enough " + symbol + " to complete sale. You currently own " + strconv.Itoa(int(u.Shares[symbol])) + " " + symbol + ".")
			}

			u.Shares[symbol] -= uint(count)
		} else {
			return errors.New("You do not own any " + symbol + ".")
		}

		err = u.Save()

		if err != nil {
			return err
		}

		_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content: "",
			Embed: &discordgo.MessageEmbed{
				Title: strconv.Itoa(count) + "x$" + symbol + " Purchase Complete",
				Fields: []*discordgo.MessageEmbedField{
					&discordgo.MessageEmbedField{
						Name:  "Buy price",
						Value: usdFormatter.FormatMoney(sharePx),
					},
					&discordgo.MessageEmbedField{
						Name:  "Total cost",
						Value: usdFormatter.FormatMoney(totalPx),
					},
				},
				Color: 0x3E606F,
			},
		})

		if err != nil {
			return err
		}

		return nil
	})
}
