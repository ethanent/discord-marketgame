package main

import (
	"encoding/csv"
	"errors"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type candidatesList struct {
	candidates [][]interface{}
}

func (c *candidatesList) Len() int {
	return len(c.candidates)
}

func (c *candidatesList) Swap(i, j int) {
	hold := c.candidates[i]
	c.candidates[i] = c.candidates[j]
	c.candidates[j] = hold
}

func (c *candidatesList) Less(i, j int) bool {
	if c.candidates[i][1].(float64) > c.candidates[j][1].(float64) {
		return true
	} else {
		return false
	}
}

func cmdPrimary(s *discordgo.Session, m *discordgo.Message, args []string) error {
	if len(args) < 1 {
		return errors.New("Missing argument.\nUsage: !primary <state name | national>")
	}

	args[0] = strings.ToLower(args[0])

	resp, err := http.Get("https://projects.fivethirtyeight.com/2020-primary-data/pres_primary_avgs_2020.csv")

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	rdr := csv.NewReader(resp.Body)

	chances := map[string]float64{}

	var onlyDate string

	for i := 0; ; i++ {
		cOdds, err := rdr.Read()

		if err != nil {
			return err
		}

		// Continue on col names row

		if i == 0 {
			continue
		}

		// Save date for first row

		if i == 1 {
			onlyDate = cOdds[2]
		}

		// If wrong date, exit.

		if cOdds[2] != onlyDate {
			// Latest date has been passed, exit parsing.
			break
		}

		if strings.ToLower(cOdds[1]) == args[0] || (args[0] == "national" && cOdds[1] == "") {
			// Candidate odds cOdds is relevant

			pctEst, err := strconv.ParseFloat(cOdds[4], 64)

			if err != nil {
				return err
			}

			chances[cOdds[3]] = pctEst
		}
	}

	// Sort

	cList := candidatesList{
		candidates: [][]interface{}{},
	}

	for cName, cOdds := range chances {
		cList.candidates = append(cList.candidates, []interface{}{cName, cOdds})
	}

	sort.Sort(&cList)

	// Prepare message

	candidateOddsFields := []*discordgo.MessageEmbedField{}

	for _, candidateOdds := range cList.candidates {
		candidateOddsFields = append(candidateOddsFields, &discordgo.MessageEmbedField{
			Name:  candidateOdds[0].(string),
			Value: strconv.FormatFloat(candidateOdds[1].(float64), 'f', 1, 64) + "%",
		})
	}

	if len(candidateOddsFields) == 0 {
		return errors.New("Bad state name.\nUsage: !primary <state name | national>")
	}

	_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content: "",
		Embed: &discordgo.MessageEmbed{
			Title:       ":ballot_box: Candidate Predictions for " + strings.ToUpper(args[0]),
			Description: "As of " + onlyDate,
			Color:       0x3E606F,
			Fields:      candidateOddsFields,
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "Source: FiveThirtyEight Average",
				IconURL: "https://fivethirtyeight.com/wp-content/themes/espn-fivethirtyeight/assets/images/fivethirtyeight-logo.png",
			},
		},
	})

	return err
}
