// Copyright (c) 2017 Alex Pliutau, Wizeline

package bot

import (
	"io/ioutil"
	"testing"
)

func init() {
	a := NewAPIMock()
	Start(a, ioutil.Discard)
}

func TestStartConversation(t *testing.T) {
	err := startConversation("testuser")
	if err != nil {
		t.Fatalf("Failed to startConversation, got: %v", err)
	}

	// Second time there are no available users
	err = startConversation("testuser")
	if err == nil {
		t.Fatalf("Error expected")
	}
	err = startConversation("teststranger")
	if err == nil {
		t.Fatalf("Error expected")
	}

	if len(conversations) != 2 {
		t.Fatalf("2 active conversation expected, got %v", conversations)
	}

	endConversation("testuser")
	endConversation("teststranger")
}

func TestEndConversation(t *testing.T) {
	err := endConversation("testuser")
	if err == nil {
		t.Fatalf("Error expected")
	}

	startConversation("testuser")
	err = endConversation("testuser")
	if err != nil {
		t.Fatalf("Failed to endConversation, got: %v", err)
	}

	startConversation("teststranger")
	err = endConversation("testuser")
	if err != nil {
		t.Fatalf("Failed to endConversation, got: %v", err)
	}

	startConversation("testuser")
	err = endConversation("inactive")
	endConversation("testuser")
	if err == nil {
		t.Fatalf("Error expected")
	}
}

func TestForwardMessage(t *testing.T) {
	startConversation("testuser")
	err := forwardMessage("testuser", "msg1")
	if err != nil {
		t.Fatalf("Failed to forwardMessage, got %v", err)
	}
	err = forwardMessage("teststranger", "msg2")
	if err != nil {
		t.Fatalf("Failed to forwardMessage, got %v", err)
	}
	err = forwardMessage("inactive", "msg3")
	if err == nil {
		t.Fatal("Error expected")
	}
	endConversation("testuser")
}

func TestGetAvailableUsers(t *testing.T) {
	users, _ := getAvailableUsers("")
	if len(users) != 2 {
		t.Fatalf("Expected 2 user, got %v", users)
	}
}

func TestFindRandomUser(t *testing.T) {
	_, randomErr := findRandomUser("")
	if randomErr != nil {
		t.Fatalf("Not expected error, got %v", randomErr)
	}
}

func TestFindRandomUserWithExclude(t *testing.T) {
	_, randomErr := findRandomUser("testuser")
	if randomErr != nil {
		t.Fatalf("Not expected error, got %v", randomErr)
	}
}

func TestGetRandomUser(t *testing.T) {
	userID := "testuser"
	random := getRandomUser([]string{userID})
	if random != userID {
		t.Fatalf("Expected %s, got %s", userID, random)
	}
}
