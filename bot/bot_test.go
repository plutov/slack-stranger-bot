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

func TestGetAvailableUsers(t *testing.T) {
	users, _ := getAvailableUsers("")
	if len(users) != 1 {
		t.Fatalf("Expected 1 user, got %v", users)
	}
}

func TestFindRandomUser(t *testing.T) {
	_, randomErr := findRandomUser("")
	if randomErr != nil {
		t.Fatalf("Not expected error, got %v", randomErr)
	}
}

func TestFindRandomUserWithExclude(t *testing.T) {
	userID := "testuser"
	_, randomErr := findRandomUser(userID)
	if randomErr == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestGetRandomUser(t *testing.T) {
	userID := "testuser"
	random := getRandomUser([]string{userID})
	if random != userID {
		t.Fatalf("Expected %s, got %s", userID, random)
	}
}
