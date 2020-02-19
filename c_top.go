package main

import (
	"errors"
	"sort"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

type SortMode int

const (
	// NetWorth sorts by net worth
	NetWorth SortMode = 0

	// Delta sorts by growth during season
	Delta SortMode = 1
)

// Leaderboard is a sortable leaderboard
type Leaderboard struct {
	users    []*User
	sortMode SortMode
}

func (l *Leaderboard) Len() int {
	return len(l.users)
}

func (l *Leaderboard) Less(i, j int) bool {
	if l.sortMode == NetWorth {
		// Sort by net worth

		iNet, err := l.users[i].NetWorth(false)

		if err != nil {
			return false
		}

		jNet, err := l.users[j].NetWorth(false)

		if err != nil {
			return false
		}

		return iNet > jNet
	} else if l.sortMode == Delta {
		// Sort by season delta

		iDelta, err := l.users[i].SeasonDelta(false)

		if err != nil {
			return false
		}

		jDelta, err := l.users[j].SeasonDelta(false)

		if err != nil {
			return false
		}

		if iDelta > jDelta {
			// User i growth > User j growth

			return true
		}

		// User j growth >= User i growth

		return false
	}

	return false
}

func (l *Leaderboard) Swap(i, j int) {
	h := l.users[i]
	l.users[i] = l.users[j]
	l.users[j] = h
}

func cmdTop(s *discordgo.Session, m *discordgo.Message, args []string) error {
	useSortMode := Delta

	if len(args) > 0 {
		switch args[0] {
		case "delta":
			useSortMode = Delta
		case "net":
			useSortMode = NetWorth
		default:
			return errors.New("Unexpected sort mode '" + args[0] + "'. Expected 'delta' or 'net'.")
		}
	}

	userIDs, err := ListUsers()

	if err != nil {
		return err
	}

	lb := Leaderboard{
		users:    []*User{},
		sortMode: useSortMode,
	}

	for _, userID := range userIDs {
		// Here one could add a check for whether or not the user is a member of the current server

		u, err := GetUser(userID)

		if err != nil {
			return err
		}

		lb.users = append(lb.users, u)
	}

	// Sort and cap users

	sort.Sort(&lb)

	if len(lb.users) > 5 {
		lb.users = lb.users[:5]
	}

	// Build message to send

	title := ""

	switch lb.sortMode {
	case NetWorth:
		title = ":archery: Net Worth Leaderboard :trophy:"
	case Delta:
		title = ":archery: Season Growth Leaderboard :trophy:"
	default:
		return errors.New("Unknown sortMode for leaderboard.")
	}

	buildEmbedLb := &discordgo.MessageEmbed{
		Title:  title,
		Fields: []*discordgo.MessageEmbedField{},
		Color:  0x46E8B2,
	}

	for place, u := range lb.users {
		val := ""

		if lb.sortMode == NetWorth {
			nw, err := u.NetWorth(false)

			if err != nil {
				return err
			}

			val = "Net worth: " + usdFormatter.FormatMoney(nw)
		} else if lb.sortMode == Delta {
			pc, err := u.SeasonDelta(false)

			if err != nil {
				return err
			}

			val = "Season growth: " + strconv.FormatFloat(pc*100, 'f', 2, 64) + "%"
		}

		dcU, err := s.User(u.ID)

		if err != nil {
			return err
		}

		buildEmbedLb.Fields = append(buildEmbedLb.Fields, &discordgo.MessageEmbedField{
			Name:  strconv.Itoa(place+1) + ". " + dcU.Username + "#" + dcU.Discriminator,
			Value: val,
		})
	}

	_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content: "",
		Embed:   buildEmbedLb,
	})

	return err
}
