package main

import (
	"errors"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func cmdList(s *discordgo.Session, m *discordgo.Message, args []string) error {
	// User for balance
	var user *User
	var username string
	
	var err error
	
	if len(args) > 0 {
		member, err := searchGuild(s, strings.Join(args, " "), m.GuildID)
		if err != nil {
			return err
		}

		if member.User.Bot {
			return errors.New("Specified user is a bot")
		}

		user, err = GetUser(member.User.ID)
		if err != nil {
			return err
		}

		username = member.User.Username
	} else {
		// Get current user's balance
		user, err = GetUser(m.Author.ID)
		if err != nil {
			return err
		}

		username = m.Author.Username
	}

	positions := []positionData{}
	var netWorth float64 = 0

	positions = append(positions, positionData{
		title:       "Cash",
		value:       user.Balance,
		dayIncrease: 0,
		percentNet:  -1,
	})

	netWorth += user.Balance

	for symbol, count := range user.Shares {
		sharePx, err := getLivePrice(symbol, false)

		if err != nil {
			return err
		}

		pxDeltaPercent, err := getDayChange(symbol)

		if err != nil {
			return err
		}

		var positionValue float64 = sharePx * float64(count)

		positions = append(positions, positionData{
			title:       strconv.Itoa(int(count)) + " x " + symbol,
			value:       positionValue,
			dayIncrease: pxDeltaPercent,
			percentNet:  -1,
		})

		netWorth += positionValue
	}

	embed := discordgo.MessageEmbed{
		Title:  username + "'s Positions",
		Fields: []*discordgo.MessageEmbedField{},
	}

	for _, pos := range positions {
		pos.percentNet = pos.value / netWorth

		showTitle := pos.title

		if pos.title != "Cash" {
			formattedDI := strconv.FormatFloat(pos.dayIncrease, 'f', 2, 64)

			if pos.dayIncrease < 0 {
				showTitle += " (" + formattedDI + "%)"
			} else {
				showTitle += " (+" + formattedDI + "%)"
			}
		}

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  showTitle,
			Value: usdFormatter.FormatMoney(pos.value) + " (" + strconv.FormatFloat(pos.percentNet*100, 'f', 1, 64) + "% of portfolio)",
		})
	}

	_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Embed: &embed,
	})

	if err != nil {
		return err
	}

	return nil
}
