// Copyright (c) 2018 Alex Pliutau

package main

import (
	"os"

	"github.com/plutov/slack-stranger-bot/bot"
)

func main() {
	api := bot.NewAPISlack(os.Getenv("SLACK_TOKEN"))
	b := bot.New(api)
	b.Start(os.Stdout)
}
