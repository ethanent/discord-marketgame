package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func cmdReset(s *discordgo.Session, m *discordgo.Message, args []string) error {
	u, err := GetUser(m.Author.ID)

	if err != nil {
		return err
	}

	allowResetAfter := u.LastReset.Add(time.Hour * 24 * 7 * 4)

	if allowResetAfter.After(time.Now()) {
		return errors.New("You must wait a while before resetting. (At least 4 weeks between resets.)")
	}

	fmt.Println("Resetting user " + u.ID)

	u.Shares = map[string]int{}
	u.Balance = config["game"].(map[string]interface{})["startBalance"].(float64)
	u.SeasonStartNW = u.Balance
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
}
