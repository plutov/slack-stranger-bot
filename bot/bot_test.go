// Copyright (c) 2017 Alex Pliutau, Wizeline

package bot

import (
	"io/ioutil"
	"testing"

	"github.com/nlopes/slack"
)

var (
	bot *Bot
)

// APIMock struct
type APIMock struct{}

// NewAPIMock contructor
func NewAPIMock() *APIMock {
	return new(APIMock)
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
			ID:       "teststranger",
			Presence: "active",
		},
		slack.User{
			ID:       "inactive",
			Presence: "inactive",
		},
		slack.User{
			ID:       "bot",
			Presence: "active",
			IsBot:    true,
		},
	}, nil
}

func (a *APIMock) postMsg(channel, text string) error {
	return nil
}

func init() {
	bot = New(NewAPIMock())
	bot.Start(ioutil.Discard)
}

func BenchmarkGetChannelAndMsgFromText(b *testing.B) {
	for n := 0; n < b.N; n++ {
		bot.getChannelIDAndMsgFromText("<#C7KC1D50C|vn-bots> <#C7KC1D50C|vn-bots> hi hello")
	}
}

func TestStartConversation(t *testing.T) {
	err := bot.startConversation("testuser")
	if err != nil {
		t.Fatalf("Failed to startConversation, got: %v", err)
	}

	// Second time there are no available users
	err = bot.startConversation("testuser")
	if err == nil {
		t.Fatalf("Error expected")
	}
	err = bot.startConversation("teststranger")
	if err == nil {
		t.Fatalf("Error expected")
	}

	if len(bot.conversations) != 2 {
		t.Fatalf("2 active conversation expected, got %v", bot.conversations)
	}

	bot.endConversation("testuser", "teststranger")
}

func TestEndConversation(t *testing.T) {
	bot.startConversation("testuser")
	bot.endConversation("testuser", "teststranger")

	if len(bot.conversations) != 0 {
		t.Fatalf("0 active conversation expected, got %v", bot.conversations)
	}
}

func TestForwardMessage(t *testing.T) {
	bot.startConversation("testuser")
	err := bot.forwardMessage("testuser", "teststranger", "msg1")
	if err != nil {
		t.Fatalf("Failed to forwardMessage, got %v", err)
	}
	bot.endConversation("testuser", "teststranger")
}

func TestGetAvailableUsers(t *testing.T) {
	users, _ := bot.getAvailableUsers("")
	if len(users) != 2 {
		t.Fatalf("Expected 2 user, got %v", users)
	}
}

func TestFindRandomUser(t *testing.T) {
	_, randomErr := bot.findRandomUser("")
	if randomErr != nil {
		t.Fatalf("Not expected error, got %v", randomErr)
	}
}

func TestFindRandomUserWithExclude(t *testing.T) {
	_, randomErr := bot.findRandomUser("testuser")
	if randomErr != nil {
		t.Fatalf("Not expected error, got %v", randomErr)
	}
}

func TestGetRandomUser(t *testing.T) {
	userID := "testuser"
	random := bot.getRandomUser([]string{userID})
	if random != userID {
		t.Fatalf("Expected %s, got %s", userID, random)
	}
}

func TestGetChannelAndMsgFromText(t *testing.T) {
	chanID, msg := bot.getChannelIDAndMsgFromText("<#C7KC1D50C|vn-bots> hi hello")
	if chanID != "C7KC1D50C" {
		t.Fatalf("Expected %s, got %s", "C7KC1D50C", chanID)
	}
	if msg != "hi hello" {
		t.Fatalf("Expected %s, got %s", "hi hello", msg)
	}
}

func TestGetChannelAndMsgFromTextWithoutChannel(t *testing.T) {
	chanID, msg := bot.getChannelIDAndMsgFromText("hi hello")
	if chanID != "" {
		t.Fatalf("Expected %s, got %s", "", chanID)
	}
	if msg != "" {
		t.Fatalf("Expected %s, got %s", "", msg)
	}
}

func TestSanitizeMsg(t *testing.T) {
	original := " hi my name is @alex.pliutau @alex.pliutau"
	clean := sanitizeMsg(original)
	if clean != "hi my name is *** ***" {
		t.Fatalf("wrong sanitized msg, got %s", clean)
	}

	original2 := "@alex"
	clean2 := sanitizeMsg(original2)
	if clean2 != "***" {
		t.Fatalf("wrong sanitized msg, got %s", clean2)
	}
}
