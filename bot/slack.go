// Copyright (c) 2017 Alex Pliutau, Wizeline

package bot

import (
	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
)

// APISlack struct
type APISlack struct {
	api *slack.Client
}

// APIMock struct
type APIMock struct{}

// IAPI interface: slack or mock
type IAPI interface {
	newRTM() *slack.RTM
	getUsers() ([]slack.User, error)
	postMsg(string, string) error
}

// NewAPISlack contructor
func NewAPISlack(token string) *APISlack {
	a := new(APISlack)
	a.api = slack.New(token)
	return a
}

// NewAPIMock contructor
func NewAPIMock() *APIMock {
	return new(APIMock)
}

func (a *APISlack) newRTM() *slack.RTM {
	return a.api.NewRTM()
}

func (a *APISlack) getUsers() ([]slack.User, error) {
	return a.api.GetUsers()
}

func (a *APISlack) postMsg(channel, text string) error {
	_, _, msgErr := a.api.PostMessage(channel, text, slack.PostMessageParameters{
		AsUser: true,
	})
	if msgErr != nil {
		log.Info("[postMessage] " + msgErr.Error())
	}

	return msgErr
}

func (a *APIMock) newRTM() *slack.RTM {
	return nil
}

func (a *APIMock) getUsers() ([]slack.User, error) {
	return []slack.User{
		slack.User{
			ID:       "testuser",
			Presence: "active",
		},
		slack.User{
			ID:       "inactive",
			Presence: "inactive",
		},
	}, nil
}

func (a *APIMock) postMsg(channel, text string) error {
	return nil
}
