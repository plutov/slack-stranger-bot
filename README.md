### Slack Stranger Bot

Meet strangers in your Slack Workspace, explore new people. Fully anonymous and secure, bot doesn't store any data.

### How it works?

 - User opens a conversation with Stranger Bot.
 - Types `hi`.
 - Stranger Bot finds random active user who doesn't participate currently in another Stranger conversation.
 - Bot will act as a proxy between two connected users. Never mentioning who sends a message.
 - Any user can type `bye` to finish the conversation, and type `hi` again to start a new random one.

### Anonymous messages to the channel

You can send private message to the Stranger Bot started with channel name to send message to the channel.

 - Add Stranger Bot to the channel.
 - Send `#channel-name message` to the Bot in private conversation.

### Run Stranger Bot

1. [Create an app](https://api.slack.com/apps/) in Slack
2. Add `chat:write` OAuth scope
3. Install in your Workspace
4. Retrieve `Bot User OAuth Token`

```bash
docker build -t slack-stranger-bot .
docker run slack-stranger-bot -e SLACK_TOKEN=<token>
```

### Run Unit Tests

```
go test -v ./...
```
