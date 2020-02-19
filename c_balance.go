package main

import "github.com/bwmarrin/discordgo"

func cmdBalance(s *discordgo.Session, m *discordgo.Message, args []string) error {
	user, err := GetUser(m.Author.ID)

	if err != nil {
		return err
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
}
