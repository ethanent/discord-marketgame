package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type positionData struct {
	title       string
	value       float64
	dayIncrease float64
	percentNet  float64
}

func registerCommands() {
	registerAlternate("h", "help")

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

	registerAlternate("bal", "balance")
	registerAlternate("b", "balance")
	registerAlternate("money", "balance")
	registerAlternate("$", "balance")

	registerCommand("balance", func(s *discordgo.Session, m *discordgo.Message, args []string) error {
		user, err := GetUser(m.Author.ID)

		if err != nil {
			return err
		}

		var stocksValue float64 = 0

		for symbol, count := range user.Shares {
			oneShareValue, err := getLivePrice(symbol)

			if err != nil {
				return err
			}

			stocksValue += oneShareValue * float64(count)
		}

		_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content: "",
			Embed: &discordgo.MessageEmbed{
				Title: m.Author.Username + "'s Account",
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
	})

	registerAlternate("stocks", "shares")
	registerAlternate("list", "shares")
	registerAlternate("positions", "shares")

	registerCommand("shares", func(s *discordgo.Session, m *discordgo.Message, args []string) error {
		user, err := GetUser(m.Author.ID)

		if err != nil {
			return err
		}

		positions := []positionData{}
		var netWorth float64 = 0

		positions = append(positions, positionData{
			title:       "Cash",
			value:       user.Balance,
			dayIncrease: 0,
			percentNet:  -1,
		})

		netWorth += user.Balance

		for symbol, count := range user.Shares {
			sharePx, err := getLivePrice(symbol)

			if err != nil {
				return err
			}

			pxDeltaPercent, err := getDayChange(symbol)

			if err != nil {
				return err
			}

			var positionValue float64 = sharePx * float64(count)

			positions = append(positions, positionData{
				title:       strconv.Itoa(int(count)) + " x " + symbol,
				value:       positionValue,
				dayIncrease: pxDeltaPercent,
				percentNet:  -1,
			})

			netWorth += positionValue
		}

		embed := discordgo.MessageEmbed{
			Title:  m.Author.Username + "'s Positions",
			Fields: []*discordgo.MessageEmbedField{},
		}

		for _, pos := range positions {
			pos.percentNet = pos.value / netWorth

			showTitle := pos.title

			if pos.title != "Cash" {
				formattedDI := strconv.FormatFloat(pos.dayIncrease, 'f', 2, 64)

				if pos.dayIncrease < 0 {
					showTitle += " (" + formattedDI + "% today)"
				} else {
					showTitle += " (+" + formattedDI + "% today)"
				}
			}

			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:  showTitle,
				Value: usdFormatter.FormatMoney(pos.value) + " (" + strconv.FormatFloat(pos.percentNet*100, 'f', 1, 64) + "% of portfolio)",
			})
		}

		_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Embed: &embed,
		})

		if err != nil {
			return err
		}

		return nil
	})

	registerAlternate("px", "price")

	registerCommand("price", func(s *discordgo.Session, m *discordgo.Message, args []string) error {
		if len(args) < 1 {
			return errors.New("Missing arguments.\nUsage: !price <symbol>")
		}

		symbol := strings.ToUpper(args[0])

		price, err := getLivePrice(symbol)

		if err != nil {
			return err
		}

		user, err := GetUser(m.Author.ID)

		if err != nil {
			return err
		}

		pxChange, err := getDayChange(symbol)

		if err != nil {
			return err
		}

		addSymbol := ""

		if pxChange > 0 {
			addSymbol = "+"
		}

		userEquity, ok := user.Shares[symbol]

		if !ok {
			userEquity = 0
		}

		_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content: "",
			Embed: &discordgo.MessageEmbed{
				Title: symbol + " Stock",
				Fields: []*discordgo.MessageEmbedField{
					&discordgo.MessageEmbedField{
						Name:  "Share price",
						Value: usdFormatter.FormatMoney(price) + " (" + addSymbol + strconv.FormatFloat(pxChange, 'f', 2, 64) + "% from yesterday)",
					},
					&discordgo.MessageEmbedField{
						Name:  m.Author.Username + "'s equity",
						Value: usdFormatter.FormatMoney(price*float64(userEquity)) + " (" + strconv.Itoa(int(userEquity)) + " shares)",
					},
				},
				Color: 0x3E606F,
			},
		})

		return err
	})

	registerCommand("reset", func(s *discordgo.Session, m *discordgo.Message, args []string) error {
		u, err := GetUser(m.Author.ID)

		if err != nil {
			return err
		}

		allowResetAfter := u.LastReset.Add(time.Hour * 24 * 7 * 4)

		if allowResetAfter.After(time.Now()) {
			return errors.New("You must wait a while before resetting. (At least 4 weeks between resets.)")
		}

		fmt.Println("Resetting user " + u.ID)

		u.Shares = map[string]uint{}
		u.Balance = config["game"].(map[string]interface{})["startBalance"].(float64)
		u.LastReset = time.Now()

		u.Save()

		_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content: "",
			Embed: &discordgo.MessageEmbed{
				Title:       ":tada: Reset Complete",
				Description: "Your cash and shares have been destroyed.",
				Color:       0x3E606F,
			},
		})

		if err != nil {
			return err
		}

		return nil
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

		if count == 0 {
			return errors.New("You must buy at least one share.")
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
			return errors.New("You can't afford to buy " + strconv.Itoa(count) + " x " + symbol)
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
				Title: ":tada: " + strconv.Itoa(count) + " x " + symbol + " Purchase Complete",
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
				Color: 0x46E8B2,
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

		if count == 0 {
			return errors.New("You must sell at least one share.")
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
				return errors.New("You do not own enough $" + symbol + " to complete sale. You currently own " + strconv.Itoa(int(u.Shares[symbol])) + " " + symbol + ".")
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
				Title: ":tada: " + strconv.Itoa(count) + " x " + symbol + " Sale Complete",
				Fields: []*discordgo.MessageEmbedField{
					&discordgo.MessageEmbedField{
						Name:  "Sell price",
						Value: usdFormatter.FormatMoney(sharePx),
					},
					&discordgo.MessageEmbedField{
						Name:  "Total received",
						Value: usdFormatter.FormatMoney(totalPx),
					},
				},
				Color: 0x46E8B2,
			},
		})

		if err != nil {
			return err
		}

		return nil
	})

}
