### Slack Stranger Bot [![Build Status](https://travis-ci.org/wizeline/slack-stranger-bot.svg?branch=master)](https://travis-ci.org/wizeline/slack-stranger-bot)

Meet strangers in your company, explore new people. Fully anonymous and secure, bot doesn't store any data.

### How it works?

 - User opens conversation with Stranger Bot.
 - Types `hi`.
 - Stranger Bot finds random active user who doesn't participate currently in Stranger conversation.
 - Bot will forward all next messages sent by user to Bot to the Stranger user. Without mentioning who sent this message.
 - Any user can type `bye` to finish the conversation, and type `hi` again to start a new random one.

### Start Bot

 - Create App in Slack and copy `token`
 - Install [Docker](https://docs.docker.com/engine/installation/)
 - Run `docker build -t stranger . && docker run stranger -e SLACK_TOKEN=<token>` with valid token

### Run Unit Tests

Install glide and dependencies:

```
curl --silent https://glide.sh/get | sh
glide install
```

```
go test ./...
```
