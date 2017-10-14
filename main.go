package main

import (
	"github.com/wizeline/slack-stranger-bot/bot"
	"os"
)

func main() {
	bot.Start(os.Getenv("SLACK_TOKEN"))
}
