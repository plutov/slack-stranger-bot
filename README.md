### Slack Stranger Bot

Meet strangers in your company, explore new people. Fully anonymous and secure, bot doesn't store any data.

### How it works?

 - User opens conversation with Stranger Bot.
 - Types `hi`.
 - Stranger Bot finds random user who doesn't participate currently in Stranger conversation.
 - All following messages sent by user to Bot will be forwarded by Bot to Stranger user.
 - Any user can type `bye` to finish conversation, and type `hi` again to start a new random one.

### Start Bot

 - Create App in Slack and copy `token`
 - Install [Docker](https://docs.docker.com/engine/installation/)
 - Run `docker build -t stranger . && docker run stranger -e SLACK_TOKEN=<token>` with valid token
