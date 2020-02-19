package main

import (
	"errors"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func cmdStop(s *discordgo.Session, m *discordgo.Message, args []string) error {
	if len(args) < 3 {
		return errors.New("Missing arguments.\nUsage: !stop <symbol> <price> <count>")
	}

	// Parse args

	count, err := strconv.Atoi(args[2])

	if err != nil {
		return err
	}

	price, err := strconv.ParseFloat(args[1], 64)

	if err != nil {
		return err
	}

	symbol := strings.ToUpper(args[0])

	// Get user

	u, err := GetUser(m.Author.ID)

	if err != nil {
		return err
	}

	// Ensure there is not currently a stop order

	_, ok := u.Stops[symbol]

	if ok == true {
		return errors.New("You have a pending stop order for " + symbol + " which blocks this stop placement. You may cancel the current order using the command \"!cancel " + symbol + "\"")
	}

	// Ensure user has enough shares to sell

	if u.Shares[symbol] < count {
		return errors.New("You do not own enough " + symbol + " to create this stop order.")
	}

	// Ensure stop price is below current market price (this stops exploits using IVT!)

	curPx, err := getLivePrice(symbol, true)

	if err != nil {
		return err
	}

	if price > curPx {
		return errors.New("The price specified (" + usdFormatter.FormatMoney(price) + ") is above market price for " + symbol + ".")
	}

	// Place stop order

	u.Stops[symbol] = StopOrder{
		Price: price,
		Count: count,
	}

	err = u.Save()

	if err != nil {
		return err
	}

	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content: "",
		Embed: &discordgo.MessageEmbed{
			Title: ":ballot_box_with_check: " + strconv.Itoa(count) + " x " + symbol + " Stop Order Placed",
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:  "Stop price",
					Value: usdFormatter.FormatMoney(price),
				},
				&discordgo.MessageEmbedField{
					Name:  "Quantity",
					Value: strconv.Itoa(count),
				},
			},
			Color: 0x46E8B2,
		},
	})

	return nil
}
