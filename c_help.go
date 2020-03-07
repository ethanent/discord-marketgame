package main

import "github.com/bwmarrin/discordgo"

func cmdHelp(s *discordgo.Session, m *discordgo.Message, args []string) error {
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
					Value:  "\n!help\n!price <ticker>\n!bal\n!list\n!top [net | delta]",
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Activity Commands",
					Value:  "\n!buy <symbol> <count>\n!sell <symbol> <count>\n!panic <symbol>\n!panicall\n!stop <symbol> <price> <count>\n!cancel <symbol>\n!reset",
					Inline: true,
				},
			},
			Color: 0x3E606F,
		},
	})

	return err
}
