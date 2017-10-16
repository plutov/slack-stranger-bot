package main

import (
	"os"

	"github.com/wizeline/slack-stranger-bot/bot"
)

func main() {
	api := bot.NewAPISlack(os.Getenv("SLACK_TOKEN"))
	bot.Start(api, os.Stdout)
}
