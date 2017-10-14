// Copyright (c) 2017 Alex Pliutau, Wizeline

package bot

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"sync"
	"testing"
)

func init() {
	mu = &sync.Mutex{}
	conversations = make(map[string]string)
	log.SetOutput(ioutil.Discard)
}

func TestGetRandomUser(t *testing.T) {
	userID := "testuser"
	random := getRandomUser([]string{userID})
	if random != userID {
		t.Fatalf("Expected %s, got %s", userID, random)
	}
}
