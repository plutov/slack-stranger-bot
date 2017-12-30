// Copyright (c) 2017 Alex Pliutau, Wizeline

package bot

import (
	"fmt"
	"io"
	"math/rand"
	"strings"
	"sync"

	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
)

// Bot type
type Bot struct {
	// initiator => stranger map
	conversations map[string]string
	// channel for all incoming messages, handled in startRTM
	pipeline chan *slack.MessageEvent
	api      IAPI
	mu       *sync.Mutex
}

const (
	startCommand   = "hi"
	endCommand     = "bye"
	connMsg        = "Connecting to a random Stranger ..."
	foundMsg       = "Stranger found! Say hello, and please be polite :wave:. It's anonymous! _Type *bye* to finish the conversation_"
	strangerMsg    = "One random Stranger just connected to you. Wanna talk? Type something here and I will forward it to Stranger anonymously. _Or type *bye* to finish the conversation_"
	byeMsg         = "Bye! You finished the conversation with the Stranger. _Type *hi* again if you want to start a new random one._"
	byeStrangerMsg = "Bye! Stranger finished the conversation with you. _Type *hi* again if you want to start a new random one._"
	notFoundMsg    = "Sorry, cannot find available online Stranger right now :disappointed:"
)

// New inits new bot
func New(api IAPI) *Bot {
	b := new(Bot)
	b.mu = &sync.Mutex{}
	b.conversations = make(map[string]string)
	b.pipeline = make(chan *slack.MessageEvent)
	b.api = api

	return b
}

// Start bot entry func
func (b *Bot) Start(logOut io.Writer) {
	log.SetOutput(logOut)
	log.Info("[main] Stranger Bot started.")

	b.startRTM()
}

func (b *Bot) startRTM() {
	rtm := b.api.newRTM()
	if rtm == nil {
		log.Error("rtm object is nil")
		return
	}

	go rtm.ManageConnection()

	go func() {
		for {
			ev := <-b.pipeline
			b.handleMessageEvent(ev)
		}
	}()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			// Do not handle bot messages
			// Do not handle non-private messages
			if ev == nil || len(ev.BotID) > 0 || !b.isPrivateMsg(ev) {
				continue
			}

			b.pipeline <- ev
		}
	}
}

func (b *Bot) handleMessageEvent(ev *slack.MessageEvent) {
	text := b.sanitizeMsg(ev.Msg.Text)

	// Send anonymous message to the channel if message starts with channel name
	chanID, msg := b.getChannelIDAndMsgFromText(text)
	if len(chanID) > 0 && len(msg) > 0 {
		b.api.postMsg(chanID, msg)
		return
	}

	b.mu.Lock()
	stranger, found := b.conversations[ev.Msg.User]
	b.mu.Unlock()

	var err error
	possibleCommand := strings.ToLower(text)
	if possibleCommand == startCommand && !found {
		err = b.startConversation(ev.Msg.User)
	} else if possibleCommand == endCommand && found {
		err = b.endConversation(ev.Msg.User, stranger)
	} else if found {
		err = b.forwardMessage(ev.Msg.User, stranger, text)
	}

	if err != nil {
		log.Error(err.Error())
	}
}

// Remove usernames, make trim
func (b *Bot) sanitizeMsg(msg string) string {
	msg = strings.TrimSpace(msg)

	return msg
}

func (b *Bot) isPrivateMsg(ev *slack.MessageEvent) bool {
	return len(ev.Channel) > 0 && string(ev.Channel[0]) == "D"
}

// Parse message text and get channel name from the beginning of the text
func (b *Bot) getChannelIDAndMsgFromText(msg string) (string, string) {
	parts := strings.Split(msg, " ")
	if len(parts) > 1 && strings.HasPrefix(parts[0], "<#") {
		chanParts := strings.Split(parts[0], "|")

		if len(chanParts) > 0 {
			r := strings.NewReplacer("<#", "", "|", "")
			chanID := r.Replace(chanParts[0])

			// Second return val is the message without channel name
			return chanID, strings.Join(parts[1:], " ")
		}
	}

	return "", ""
}

func (b *Bot) startConversation(msgUser string) error {
	b.api.postMsg(msgUser, connMsg)

	stranger, findErr := b.findRandomUser(msgUser)

	if findErr == nil && len(stranger) > 0 {
		b.mu.Lock()
		b.conversations[msgUser] = stranger
		b.conversations[stranger] = msgUser
		b.mu.Unlock()

		// Notify current user that we found Stranger
		b.api.postMsg(msgUser, foundMsg)
		// Notify Stranger
		b.api.postMsg(stranger, strangerMsg)

		log.Info("[startConversation] ok")
		return nil
	}

	// Notify current user that we cannot find a Stranger
	b.api.postMsg(msgUser, notFoundMsg)

	return fmt.Errorf("[startConversation] stranger not found, %v", findErr)
}

// user -> bot -> user. Secure
func (b *Bot) forwardMessage(msgUser string, stranger string, text string) error {
	b.api.postMsg(stranger, text)
	log.Info("[forwardMessage] ok")
	return nil
}

func (b *Bot) endConversation(msgUser string, stranger string) error {
	// Notify Initiator and Stranger that conversation is finished
	b.api.postMsg(msgUser, byeMsg)
	b.api.postMsg(stranger, byeStrangerMsg)

	b.mu.Lock()
	delete(b.conversations, msgUser)
	delete(b.conversations, stranger)
	b.mu.Unlock()

	log.Info("[endConversation] ok")

	return nil
}

// Get all available users from Slack once
func (b *Bot) getAvailableUsers(exclude string) ([]string, error) {
	users := []string{}

	slackUsers, err := b.api.getUsers()
	if err != nil {
		return users, fmt.Errorf("[getAvailableUsers] %v", err.Error())
	}

	for _, u := range slackUsers {
		_, inConversation := b.conversations[u.ID]
		if !u.IsBot && !inConversation && u.Presence == "active" && u.ID != exclude {
			users = append(users, u.ID)
		}
	}

	return users, nil
}

func (b *Bot) findRandomUser(initiator string) (string, error) {
	availableUsers, err := b.getAvailableUsers(initiator)
	if err != nil {
		return "", err
	}
	if len(availableUsers) == 0 {
		return "", fmt.Errorf("[findRandomUser] no available users found")
	}

	randomUser := b.getRandomUser(availableUsers)

	return randomUser, nil
}

func (b *Bot) getRandomUser(list []string) string {
	return list[rand.Intn(len(list))]
}
