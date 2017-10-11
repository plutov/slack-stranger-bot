// Copyright (c) 2017 Alex Pliutau, Wizeline

package main

import (
	"github.com/nlopes/slack"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
)

type user struct {
	stranger *string
}

var (
	users map[string]*user
	api   *slack.Client
	mu    *sync.Mutex
)

const (
	startCommand   = "hi"
	endCommand     = "bye"
	connMsg        = "Connecting to a random Stranger ..."
	foundMsg       = "Stranger found! Say hello, and please be polite :wave: _Type *bye* to finish the conversation_"
	strangerMsg    = "One random Stranger just connected to you. Wanna talk? Type something here and I will forward it to Stranger anonymously. _Or type *bye* to finish the conversation_"
	byeMsg         = "Bye! You finished conversation with the Stranger. _Type *hi* again if you want to start a new random one._"
	byeStrangerMsg = "Bye! Stranger finished conversation with you. _Type *hi* again if you want to start a new random one._"
	notFoundMsg    = "Sorry, cannot find available online Stranger right now :disappointed:"
)

func main() {
	mu = &sync.Mutex{}

	api = slack.New(os.Getenv("SLACK_TOKEN"))

	getUsers()

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			mu.Lock()
			user, found := users[ev.Msg.User]
			mu.Unlock()
			possibleCommand := strings.TrimSpace(strings.ToLower(ev.Msg.Text))
			if possibleCommand == startCommand && found && user.stranger == nil {
				go startConversation(ev)
			} else if possibleCommand == endCommand && found && user.stranger != nil {
				go endConversation(ev)
			} else if found && user.stranger != nil {
				go forwardMessage(ev)
			}

		}
	}
}

// Get all users from Slack once
func getUsers() {
	users = make(map[string]*user)

	slackUsers, err := api.GetUsers()
	if err != nil {
		log.Fatal("[getUsers] " + err.Error())
	}

	for _, u := range slackUsers {
		if !u.IsBot {
			users[u.ID] = &user{}
		}
	}
}

func startConversation(ev *slack.MessageEvent) {
	params := slack.PostMessageParameters{
		AsUser: true,
	}
	postMsg(ev.Msg.User, connMsg, params)

	mu.Lock()
	stranger := findRandomUser(ev.Msg.User)
	mu.Unlock()
	if len(stranger) > 0 {
		// Notify current user that we found Stranger
		postMsg(ev.Msg.User, foundMsg, params)
		// Notify Stranger
		postMsg(stranger, strangerMsg, params)

		log.Println("Conversation started")
	} else {
		// Notify current user that we cannot find a Stranger
		postMsg(ev.Msg.User, notFoundMsg, params)

		log.Println("Cannot find stranger")
	}
}

// user -> bot -> user. Secure
func forwardMessage(ev *slack.MessageEvent) {
	params := slack.PostMessageParameters{
		AsUser: true,
	}

	mu.Lock()
	sender, found := users[ev.Msg.User]
	mu.Unlock()

	if found {
		postMsg(*sender.stranger, ev.Msg.Text, params)
	}

	// Do not log any data, just event
	log.Println("Message forwarded")
}

func endConversation(ev *slack.MessageEvent) {
	params := slack.PostMessageParameters{
		AsUser: true,
	}

	mu.Lock()
	initiator, found := users[ev.Msg.User]
	mu.Unlock()

	if found {
		postMsg(ev.Msg.User, byeMsg, params)

		mu.Lock()
		strangerID := *users[ev.Msg.User].stranger
		stranger, ok := users[strangerID]
		mu.Unlock()

		if ok {
			// Notify Stranger that conversation is finished
			postMsg(strangerID, byeStrangerMsg, params)
			stranger.stranger = nil
		}
		initiator.stranger = nil
	}

	// Do not log any data, just event
	log.Println("Conversation ended")
}

func findRandomUser(initiator string) string {
	var attemptsLeft = 25

	for attemptsLeft > 0 {
		randomID, randomUser := getRandomUser(users)
		if randomUser != nil && randomID != initiator && randomUser.stranger == nil {
			presense, err := api.GetUserPresence(randomID)
			if err != nil {
				log.Println("[findRandomUser]: " + err.Error())
			} else if presense != nil && presense.Presence == "active" {
				randomUser.stranger = &initiator
				users[initiator].stranger = &randomID
				return randomID
			}
		}
		attemptsLeft--
	}

	return ""
}

func getRandomUser(m map[string]*user) (string, *user) {
	i := rand.Intn(len(m))
	for id, u := range m {
		if i == 0 {
			return id, u
		}
		i--
	}

	return "", nil
}

func postMsg(channel, text string, params slack.PostMessageParameters) {
	_, _, msgErr := api.PostMessage(channel, text, params)
	if msgErr != nil {
		log.Println("[postMessage] " + msgErr.Error())
	}
}
