package main

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func cmdSell(s *discordgo.Session, m *discordgo.Message, args []string) error {
	if len(args) < 2 {
		return errors.New("Missing arguments.\nUsage: !sell <symbol> <count>")
	}

	count, err := strconv.Atoi(args[1])
	symbol := strings.ToUpper(args[0])

	if err != nil {
		return errors.New("Non-numeric value for count.\nUsage: !buy <symbol> <count>")
	}

	if count < 0 {
		return errors.New("You can't sell negative shares.")
	}

	if count == 0 {
		return errors.New("You must sell at least one share.")
	}

	sharePx, err := getLivePrice(symbol, true)

	if err != nil {
		return err
	}

	u, err := GetUser(m.Author.ID)

	if err != nil {
		return err
	}

	// Ensure market is open if not crypto

	if !isCrypto(symbol) {
		t, err := getLastUpdated(symbol)

		if err != nil {
			return err
		}

		if time.Since(*t).Minutes() > 8 {
			return errors.New("The market has closed and the trade cannot be completed.")
		}
	}

	// Ensure no active stops prevent sale

	_, ok := u.Stops[symbol]

	if ok == true {
		return errors.New("You have a pending stop order for " + symbol + " which blocks this sale.")
	}

	totalPx := float64(count) * sharePx

	u.Balance += totalPx

	_, ok = u.Shares[symbol]

	if ok {
		if u.Shares[symbol] < count {
			return errors.New("You do not own enough " + symbol + " to complete sale. You currently own " + strconv.Itoa(int(u.Shares[symbol])) + " " + symbol + ".")
		}

		u.Shares[symbol] -= count

		if u.Shares[symbol] <= 0 {
			delete(u.Shares, symbol)
		}
	} else {
		return errors.New("You do not own any " + symbol + ".")
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
			Thumbnail: thumbnail,
			Color:     0x46E8B2,
		},
	})

	if err != nil {
		return err
	}

	return nil
}
