package main

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

func cmdNewSeason(s *discordgo.Session, m *discordgo.Message, args []string) error {
	if m.Author.ID != config["adminUserID"].(string) {
		return errors.New("You are not authorized to start a new season.")
	}

	userIDs, err := ListUsers()

	if err != nil {
		return err
	}

	for _, userID := range userIDs {
		u, err := GetUser(userID)

		if err != nil {
			return err
		}

		nw, err := u.NetWorth(false)

		if err != nil {
			return err
		}

		u.SeasonStartNW = nw

		err = u.Save()

		if err != nil {
			return err
		}
	}

	_, err = s.ChannelMessageSend(m.ChannelID, ":white_check_mark: New season has been started.")

	return err
}
