package main

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func cmdPanic(s *discordgo.Session, m *discordgo.Message, args []string) error {
	if len(args) < 1 {
			return errors.New("Too few arguments.\nUsage: !panic <symbol>")
	}
	symbol := strings.ToUpper(args[0])

	sharePx, e := getLivePrice(symbol, true)
	if e != nil {
		return e
	}

	u, e := GetUser(m.Author.ID)
	if e != nil {
		return e
	}

	if !isCrypto(symbol) {
		t, e := getLastUpdated(symbol)

		if e != nil {
			return e
		}

		if time.Since(*t).Minutes() > 8 {
			return errors.New("The market has closed and the trade cannot be completed.")
		}
	}

	shares, ok := u.Shares[symbol]
	totalPx := float64(shares) * sharePx

	if !ok {
		return errors.New("You don't own any shares of "+symbol+".")
	}
	u.Balance += totalPx
	delete(u.Shares, symbol)

	e = u.Save()
	if e != nil {
		return e
	}

	coImageUrl, e := getLogo(symbol)
	var thumbnail *discordgo.MessageEmbedThumbnail = nil
	if e == nil {
		thumbnail = &discordgo.MessageEmbedThumbnail {
			URL: coImageUrl,
		}
	}

	_, e = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend {
		Content: "",
		Embed: &discordgo.MessageEmbed {
			Title: ":tada: Panic " + symbol + " Sale Complete",
			Fields: []*discordgo.MessageEmbedField {
				&discordgo.MessageEmbedField {
					Name:  "Amount sold",
					Value: strconv.Itoa(shares),
				},
				&discordgo.MessageEmbedField {
					Name:  "Sell price",
					Value: usdFormatter.FormatMoney(sharePx),
				},
				&discordgo.MessageEmbedField {
					Name:  "Total received",
					Value: usdFormatter.FormatMoney(totalPx),
				},
			},
			Thumbnail: thumbnail,
			Color:     0x46E8B2,
			},
		})

		if e != nil {
			return e
		}

		return nil
}
