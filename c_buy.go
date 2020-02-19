package main

import (
	"errors"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func cmdBuy(s *discordgo.Session, m *discordgo.Message, args []string) error {
	if len(args) < 2 {
		return errors.New("Missing arguments.\nUsage: !buy <symbol> <count>")
	}

	count, err := strconv.Atoi(args[1])
	symbol := strings.ToUpper(args[0])

	if err != nil {
		return err
	}

	if count < 0 {
		return errors.New("You can't buy negative shares.")
	}

	if count == 0 {
		return errors.New("You must buy at least one share.")
	}

	sharePx, err := getLivePrice(symbol, true)

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
		u.Shares[symbol] += count
	} else {
		u.Shares[symbol] = count
	}

	err = u.Save()

	if err != nil {
		return err
	}

	coImageURL, err := getLogo(symbol)

	var thumbnail *discordgo.MessageEmbedThumbnail = nil

	if err == nil {
		thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: coImageURL,
		}
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
			Thumbnail: thumbnail,
			Color:     0x46E8B2,
		},
	})

	if err != nil {
		return err
	}

	return nil
}
