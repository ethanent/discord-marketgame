package main

import (
	"errors"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func cmdPanic(s *discordgo.Session, m *discordgo.Message, args []string) error {
	if len(args) < 1 {
		return errors.New("Too few arguments.\nUsage: !panic <symbol>")
	}
	symbol := strings.ToUpper(args[0])
	
	u, e := GetUser(m.Author.ID)
	if e != nil {
		return e
	}
	
	shares, ok := u.Shares[symbol]
	
	if !ok {
		return errors.New("You don't own any shares of "+symbol+".")
	}

	return cmdSell(s, m, []string{symbol, strconv.Itoa(shares)})
}
