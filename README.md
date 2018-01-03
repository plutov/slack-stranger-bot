### Slack Stranger Bot [![Build Status](https://travis-ci.org/wizeline/slack-stranger-bot.svg?branch=master)](https://travis-ci.org/wizeline/slack-stranger-bot)

Meet strangers in your company, explore new people. Fully anonymous and secure, bot doesn't store any data.

### How it works?

 - User opens conversation with Stranger Bot.
 - Types `hi`.
 - Stranger Bot finds random active user who doesn't participate currently in Stranger conversation.
 - Bot will forward all next messages sent by user to Bot to the Stranger user. Without mentioning who sent this message.
 - Any user can type `bye` to finish the conversation, and type `hi` again to start a new random one.

### Anonymous messages to the channel

You can send private message to the Bot started with channel name to send message to the channel.

 - Add bot to the channel.
 - Send `#channel-name message` to the Bot in private conversation.

### Run Stranger

 - Create App in Slack and copy `token`
 - Install [Docker](https://docs.docker.com/engine/installation/)
 - Build Bot `docker build -t stranger .`
 - Run Bot `docker run stranger -e SLACK_TOKEN=<token>`

### Run Unit Tests

```
make test
```
