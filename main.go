// Copyright (c) 2017 Alex Pliutau, Wizeline

package main

import (
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"

	"github.com/nlopes/slack"
)

type user struct {
	id       string
	stranger *string
}

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

func main() {
	mu = &sync.Mutex{}
	conversations = make(map[string]string)

	api = slack.New(os.Getenv("SLACK_TOKEN"))

	log.Println("[main] Stranger Bot started.")

	startRTM()
}

func startRTM() {
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
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
func getAvailableUsers(exclude string) []*user {
	users := []*user{}

	slackUsers, err := api.GetUsers()
	if err != nil {
		log.Fatal("[getUsers] " + err.Error())
	}

	for _, u := range slackUsers {
		if !u.IsBot && u.Presence == "active" && u.ID != exclude {
			users = append(users, &user{
				id: u.ID,
			})
		}
	}

	return users
}

func startConversation(ev *slack.MessageEvent) {
	params := slack.PostMessageParameters{
		AsUser: true,
	}
	postMsg(ev.Msg.User, connMsg, params)

	stranger := findRandomUser(ev.Msg.User)

	if len(stranger) > 0 {
		mu.Lock()
		conversations[ev.Msg.User] = stranger
		conversations[stranger] = ev.Msg.User
		mu.Unlock()

		// Notify current user that we found Stranger
		postMsg(ev.Msg.User, foundMsg, params)
		// Notify Stranger
		postMsg(stranger, strangerMsg, params)

		log.Println("[startConversation] ok: " + ev.Msg.User + " -> " + stranger)
	} else {
		// Notify current user that we cannot find a Stranger
		postMsg(ev.Msg.User, notFoundMsg, params)

		log.Println("[startConversation] not found")
	}
}

// user -> bot -> user. Secure
func forwardMessage(ev *slack.MessageEvent) {
	params := slack.PostMessageParameters{
		AsUser: true,
	}

	mu.Lock()
	stranger, found := conversations[ev.Msg.User]
	mu.Unlock()

	if found {
		postMsg(stranger, ev.Msg.Text, params)
		log.Println("[forwardMessage] ok: " + ev.Msg.User + " -> " + stranger)
	} else {
		log.Println("[forwardMessage] unable to find stranger for " + ev.Msg.User)
	}
}

func endConversation(ev *slack.MessageEvent) {
	params := slack.PostMessageParameters{
		AsUser: true,
	}

	mu.Lock()
	stranger, found := conversations[ev.Msg.User]
	mu.Unlock()

	if found {
		// Notify Initiator and Stranger that conversation is finished
		postMsg(ev.Msg.User, byeMsg, params)
		postMsg(stranger, byeStrangerMsg, params)

		log.Println("[endConversation] ok: " + ev.Msg.User + " & " + stranger)

		mu.Lock()
		delete(conversations, ev.Msg.User)
		delete(conversations, stranger)
		mu.Unlock()
	} else {
		log.Println("[endConversation] unable to find stranger for " + ev.Msg.User)
	}
}

func findRandomUser(initiator string) string {
	var attemptsLeft = 5

	for attemptsLeft > 0 {
		availableUsers := getAvailableUsers(initiator)
		randomUser := getRandomUser(availableUsers)

		if randomUser != nil {
			return randomUser.id
		}
		attemptsLeft--
	}

	return ""
}

func getRandomUser(list []*user) *user {
	return list[rand.Intn(len(list))]
}

func postMsg(channel, text string, params slack.PostMessageParameters) {
	_, _, msgErr := api.PostMessage(channel, text, params)
	if msgErr != nil {
		log.Println("[postMessage] " + msgErr.Error())
	}
}
