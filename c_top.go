package main

import (
	"sort"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

// Leaderboard is a sortable leaderboard
type Leaderboard struct {
	users []*User
}

func (l *Leaderboard) Len() int {
	return len(l.users)
}

func (l *Leaderboard) Less(i, j int) bool {
	iNet, err := l.users[i].NetWorth(false)

	if err != nil {
		return false
	}

	jNet, err := l.users[j].NetWorth(false)

	if err != nil {
		return false
	}

	return iNet > jNet
}

func (l *Leaderboard) Swap(i, j int) {
	h := l.users[i]
	l.users[i] = l.users[j]
	l.users[j] = h
}

func cmdTop(s *discordgo.Session, m *discordgo.Message, args []string) error {
	userIDs, err := ListUsers()

	if err != nil {
		return err
	}

	lb := Leaderboard{
		users: []*User{},
	}

	for _, userID := range userIDs {
		u, err := GetUser(userID)

		if err != nil {
			return err
		}

		lb.users = append(lb.users, u)
	}

	sort.Sort(&lb)

	// Build message to send

	buildEmbedLb := &discordgo.MessageEmbed{
		Title:  ":archery: Net Worth Leaderboard :trophy:",
		Fields: []*discordgo.MessageEmbedField{},
		Color:  0x46E8B2,
	}

	for place, u := range lb.users {
		nw, err := u.NetWorth(false)

		if err != nil {
			return err
		}

		dcU, err := s.User(u.ID)

		if err != nil {
			return err
		}

		buildEmbedLb.Fields = append(buildEmbedLb.Fields, &discordgo.MessageEmbedField{
			Name:  strconv.Itoa(place+1) + ". " + dcU.Username + "#" + dcU.Discriminator,
			Value: "Net worth: " + usdFormatter.FormatMoney(nw),
		})
	}

	_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content: "",
		Embed:   buildEmbedLb,
	})

	return err
}
