package main

import (
	"coin-bot/service"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	url       string
	botPrefix string
	botKey    string
)

func main() {
	app := service.NewBot(botPrefix, botKey, url)
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
	<-make(chan struct{})
	return
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	url = os.Getenv("RSS_SERVICE")
	botPrefix = os.Getenv("BOT_PREFIX")
	botKey = os.Getenv("BOT_KEY")
}
