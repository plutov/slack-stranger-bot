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

var (
	// initiator => stranger map
	conversations map[string]string
	api           IAPI
	mu            *sync.Mutex
)

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

// Start bot entry func
func Start(a IAPI, logOut io.Writer) {
	mu = &sync.Mutex{}
	conversations = make(map[string]string)

	api = a

	log.SetOutput(logOut)
	log.Info("[main] Stranger Bot started.")

	startRTM()
}

func startRTM() {
	rtm := api.newRTM()
	if rtm == nil {
		return
	}

	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			// Do not handle bot messages
			if len(ev.BotID) > 0 {
				continue
			}

			go handleMessageEvent(ev)
		}
	}
}

func handleMessageEvent(ev *slack.MessageEvent) {
	mu.Lock()
	_, found := conversations[ev.Msg.User]
	mu.Unlock()

	var err error
	possibleCommand := strings.TrimSpace(strings.ToLower(ev.Msg.Text))
	if possibleCommand == startCommand && !found {
		err = startConversation(ev.Msg.User)
	} else if possibleCommand == endCommand && found {
		err = endConversation(ev.Msg.User)
	} else if found {
		err = forwardMessage(ev.Msg.User, ev.Msg.Text)
	}

	if err != nil {
		log.Error(err.Error())
	}
}

func startConversation(msgUser string) error {
	api.postMsg(msgUser, connMsg)

	stranger, findErr := findRandomUser(msgUser)

	if findErr == nil && len(stranger) > 0 {
		mu.Lock()
		conversations[msgUser] = stranger
		conversations[stranger] = msgUser
		mu.Unlock()

		// Notify current user that we found Stranger
		api.postMsg(msgUser, foundMsg)
		// Notify Stranger
		api.postMsg(stranger, strangerMsg)

		log.Info("[startConversation] ok")
		return nil
	}

	// Notify current user that we cannot find a Stranger
	api.postMsg(msgUser, notFoundMsg)

	return fmt.Errorf("[startConversation] stranger not found")
}

// user -> bot -> user. Secure
func forwardMessage(msgUser string, text string) error {
	mu.Lock()
	stranger, found := conversations[msgUser]
	mu.Unlock()

	if found {
		api.postMsg(stranger, text)
		log.Info("[forwardMessage] ok")
		return nil
	}

	return fmt.Errorf("[forwardMessage] unable to find stranger")
}

func endConversation(msgUser string) error {
	mu.Lock()
	stranger, found := conversations[msgUser]
	mu.Unlock()

	if found {
		// Notify Initiator and Stranger that conversation is finished
		api.postMsg(msgUser, byeMsg)
		api.postMsg(stranger, byeStrangerMsg)

		log.Info("[endConversation] ok")

		mu.Lock()
		delete(conversations, msgUser)
		delete(conversations, stranger)
		mu.Unlock()

		return nil
	}

	return fmt.Errorf("[endConversation] unable to find stranger")
}

// Get all available users from Slack once
func getAvailableUsers(exclude string) ([]string, error) {
	users := []string{}

	slackUsers, err := api.getUsers()
	if err != nil {
		return users, fmt.Errorf("[getAvailableUsers] %v", err.Error())
	}

	for _, u := range slackUsers {
		_, inConversation := conversations[u.ID]
		if !u.IsBot && !inConversation && u.Presence == "active" && u.ID != exclude {
			users = append(users, u.ID)
		}
	}

	return users, nil
}

func findRandomUser(initiator string) (string, error) {
	availableUsers, err := getAvailableUsers(initiator)
	if err != nil {
		return "", err
	}
	if len(availableUsers) == 0 {
		return "", fmt.Errorf("[findRandomUser] no available users found")
	}

	randomUser := getRandomUser(availableUsers)

	return randomUser, nil
}

func getRandomUser(list []string) string {
	return list[rand.Intn(len(list))]
}
