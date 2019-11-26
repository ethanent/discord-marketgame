package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/bwmarrin/discordgo"
	av "github.com/cmckee-dev/go-alpha-vantage"
	"github.com/leekchan/accounting"
)

var usdFormatter accounting.Accounting
var avClient *av.Client

var config map[string]interface{}

func loadConfig() {
	configPath, err := filepath.Abs("config.json")

	if err != nil {
		panic(err)
	}

	fmt.Println("Loading config from", configPath)

	configFile, err := os.Open(configPath)
	defer configFile.Close()

	if err != nil {
		panic(err)
	}

	configData, err := ioutil.ReadAll(configFile)

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(configData, &config)

	if err != nil {
		panic(err)
	}
}

func main() {
	loadConfig()
	registerCommands()

	avClient = av.NewClient(config["alphaAPIKey"].(string))

	usdFormatter = accounting.Accounting{
		Symbol:    "$",
		Precision: 2,
	}

	session, err := discordgo.New("Bot " + config["discordToken"].(string))

	if err != nil {
		panic(err)
	}

	session.AddHandler(handleMessage)

	err = session.Open()

	if err != nil {
		panic(err)
	}

	fmt.Println("Session has been opened.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	<-sc
}
