package main

import (
	"errors"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func cmdPrice(s *discordgo.Session, m *discordgo.Message, args []string) error {
	if len(args) < 1 {
		return errors.New("Missing arguments.\nUsage: !price <symbol>")
	}

	symbol := strings.ToUpper(args[0])

	price, err := getLivePrice(symbol, false)

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

	coImageURL, err := getLogo(symbol)

	var thumbnail *discordgo.MessageEmbedThumbnail = nil

	if err == nil {
		thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: coImageURL,
		}
	}

	company, err := getCompany(symbol)

	if err != nil {
		return err
	}

	_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content: "",
		Embed: &discordgo.MessageEmbed{
			Title: company.Name + " Stock (" + symbol + ")",
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:  "Share price",
					Value: usdFormatter.FormatMoney(price) + " (" + addSymbol + strconv.FormatFloat(pxChange, 'f', 2, 64) + "% from yesterday)",
				},
				&discordgo.MessageEmbedField{
					Name:  m.Author.Username + "'s equity",
					Value: usdFormatter.FormatMoney(price*float64(userEquity)) + " (" + strconv.Itoa(int(userEquity)) + " shares)",
				},
				&discordgo.MessageEmbedField{
					Name:  "Corporate profile",
					Value: "Location: " + company.City + ", " + company.State + ", " + company.Country + "\nEmployees: " + strconv.Itoa(company.Employees) + "\n[Website](" + company.Website + ")",
				},
			},
			Thumbnail: thumbnail,
			Color:     0x3E606F,
		},
	})

	return err
}
