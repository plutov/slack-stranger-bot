// Copyright (c) 2017 Alex Pliutau, Wizeline

package bot

import (
	"math/rand"
	"os"
	"strings"
	"sync"

	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
)

var (
	// initiator => stranger map
	conversations map[string]string
	api           *slack.Client
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
func Start(slackToken string) {
	mu = &sync.Mutex{}
	conversations = make(map[string]string)

	api = slack.New(slackToken)

	log.SetOutput(os.Stdout)
	log.Info("[main] Stranger Bot started.")

	startRTM()
}

func startRTM() {
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			// Do not handle bot messages
			if len(ev.BotID) > 0 {
				continue
			}

			mu.Lock()
			_, found := conversations[ev.Msg.User]
			mu.Unlock()

			possibleCommand := strings.TrimSpace(strings.ToLower(ev.Msg.Text))
			if possibleCommand == startCommand && !found {
				go startConversation(ev)
			} else if possibleCommand == endCommand && found {
				go endConversation(ev)
			} else if found {
				go forwardMessage(ev)
			}

		}
	}
}

// Get all available users from Slack once
func getAvailableUsers(exclude string) []string {
	users := []string{}

	slackUsers, err := api.GetUsers()
	if err != nil {
		log.Fatal("[getUsers] " + err.Error())
	}

	for _, u := range slackUsers {
		_, inConversation := conversations[u.ID]
		if !u.IsBot && !inConversation && u.Presence == "active" && u.ID != exclude {
			users = append(users, u.ID)
		}
	}

	return users
}

func startConversation(ev *slack.MessageEvent) {
	postMsg(ev.Msg.User, connMsg)

	stranger := findRandomUser(ev.Msg.User)

	if len(stranger) > 0 {
		mu.Lock()
		conversations[ev.Msg.User] = stranger
		conversations[stranger] = ev.Msg.User
		mu.Unlock()

		// Notify current user that we found Stranger
		postMsg(ev.Msg.User, foundMsg)
		// Notify Stranger
		postMsg(stranger, strangerMsg)

		log.Info("[startConversation] ok")
	} else {
		// Notify current user that we cannot find a Stranger
		postMsg(ev.Msg.User, notFoundMsg)

		log.Info("[startConversation] stranger not found")
	}
}

// user -> bot -> user. Secure
func forwardMessage(ev *slack.MessageEvent) {
	mu.Lock()
	stranger, found := conversations[ev.Msg.User]
	mu.Unlock()

	if found {
		postMsg(stranger, ev.Msg.Text)
		log.Info("[forwardMessage] ok")
	} else {
		log.Info("[forwardMessage] unable to find stranger")
	}
}

func endConversation(ev *slack.MessageEvent) {
	mu.Lock()
	stranger, found := conversations[ev.Msg.User]
	mu.Unlock()

	if found {
		// Notify Initiator and Stranger that conversation is finished
		postMsg(ev.Msg.User, byeMsg)
		postMsg(stranger, byeStrangerMsg)

		log.Info("[endConversation] ok")

		mu.Lock()
		delete(conversations, ev.Msg.User)
		delete(conversations, stranger)
		mu.Unlock()
	} else {
		log.Info("[endConversation] unable to find stranger")
	}
}

func findRandomUser(initiator string) string {
	availableUsers := getAvailableUsers(initiator)
	randomUser := getRandomUser(availableUsers)

	return randomUser
}

func getRandomUser(list []string) string {
	return list[rand.Intn(len(list))]
}

func postMsg(channel, text string) {
	_, _, msgErr := api.PostMessage(channel, text, slack.PostMessageParameters{
		AsUser: true,
	})
	if msgErr != nil {
		log.Info("[postMessage] " + msgErr.Error())
	}
}
