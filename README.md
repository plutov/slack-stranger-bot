### Slack Stranger Bot

Meet strangers in your company, explore new people. Fully anonymous and secure, bot doesn't store any data.

### How it works?

 - User opens conversation with Stranger Bot.
 - Types `hi`.
 - Stranger Bot finds random active user who doesn't participate currently in Stranger conversation.
 - Bot will forward all next messages sent by the user to Stranger Bot to the user. Without mentioning who sent this message.
 - Any user can type `bye` to finish the conversation, and type `hi` again to start a new random one.

### Anonymous messages to the channel

You can send private message to the Stranger Bot started with channel name to send message to the channel.

 - Add Stranger Bot to the channel.
 - Send `#channel-name message` to the Bot in private conversation.

### Run Stranger

 - Create Slack App and get API token
 - Install [Docker](https://docs.docker.com/engine/installation/)
 - Build Bot `docker build -t stranger .`
 - Run Bot `docker run stranger -e SLACK_TOKEN=<token>`

### Run Unit Tests

```
make test
```
