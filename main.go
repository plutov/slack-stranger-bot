package main

import (
	"os"

	"github.com/wizeline/slack-stranger-bot/bot"
)

func main() {
	api := bot.NewAPISlack(os.Getenv("SLACK_TOKEN"))
	b := bot.New(api)
	b.Start(os.Stdout)
}
