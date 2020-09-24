package main

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"regexp"
)

func searchGuild(s *discordgo.Session, m *discordgo.Message, query string) (*discordgo.Member, error) {
	// Check if argument is a ping
	if len(m.Mentions) == 1 {
		member, err := s.GuildMember(m.GuildID, m.Mentions[0].ID)
		if err == nil {
			return member, nil
		}
	} else if len(m.Mentions) > 1 {
		return nil, errors.New("Multiple pings in error")
	}

	// Search for user in guild
	// Using RegEx for now due to built-in nature and ease-of-use, may use something else later
	bestFit := ""
	regex, err := regexp.Compile(query)
	if err != nil {
		return nil, errors.New("Failed to compile RegEx expression")
	}
	uids, err := ListUsers()
	if err != nil {
		return nil, err
	}
	for _, uid := range uids {
		user, err := s.User(uid)
		if err != nil {
			return nil, err
		}
		match := regex.Find([]byte(user.Username))
		if match != nil {
			// Found match
			if bestFit != "" {
				// 2 users matching query
				return nil, errors.New("Couldn't find a unique user")
			}
			bestFit = uid
		}
	}
	if bestFit == "" {
		// No users found
		return nil, errors.New("Couldn't find a unique user")
	}

	return s.GuildMember(m.GuildID, bestFit)
}
