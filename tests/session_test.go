package main

import (
	"testing"
)

func TestGenerateSessionID(t *testing.T) {
	id1 := generateSessionID()
	id2 := generateSessionID()

	if len(id1) != 64 {
		t.Errorf("Expected session ID length 64, got %d", len(id1))
	}

	if id1 == id2 {
		t.Error("generateSessionID produced duplicate IDs")
	}
}

func TestCreateSession(t *testing.T) {
	username := "testuser"
	sessionID := createSession(username)

	if sessionID == "" {
		t.Fatal("createSession returned empty session ID")
	}

	if len(sessionID) != 64 {
		t.Errorf("Session ID should be 64 chars, got %d", len(sessionID))
	}
}

func TestGetSession(t *testing.T) {
	username := "testuser"
	sessionID := createSession(username)

	retrievedUsername := getSession(sessionID)
	if retrievedUsername != username {
		t.Errorf("Expected username %s, got %s", username, retrievedUsername)
	}

	// Test invalid session
	invalidUsername := getSession("invalidsessionid")
	if invalidUsername != "" {
		t.Error("getSession should return empty string for invalid session")
	}
}

func TestDeleteSession(t *testing.T) {
	username := "testuser"
	sessionID := createSession(username)

	// Verify session exists
	if getSession(sessionID) != username {
		t.Fatal("Session was not created properly")
	}

	// Delete session
	deleteSession(sessionID)

	// Verify session is gone
	retrievedUsername := getSession(sessionID)
	if retrievedUsername != "" {
		t.Error("Session was not properly deleted")
	}
}

func TestSessionConcurrency(t *testing.T) {
	// Test that concurrent session operations don't cause race conditions
	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go func(id int) {
			username := "user" + string(rune(id))
			sessionID := createSession(username)
			retrieved := getSession(sessionID)
			if retrieved != username {
				t.Errorf("Concurrent session test failed for %s", username)
			}
			deleteSession(sessionID)
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}
