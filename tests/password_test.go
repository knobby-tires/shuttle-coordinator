package main

import (
	"strings"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword123"
	hash := hashPassword(password)

	if len(hash) != 60 {
		t.Errorf("Expected hash length 60, got %d", len(hash))
	}

	if !strings.HasPrefix(hash, "$2a$") {
		t.Errorf("Hash doesn't start with bcrypt identifier: %s", hash)
	}
}

func TestHashPasswordUniqueness(t *testing.T) {
	password := "samepassword"
	hash1 := hashPassword(password)
	hash2 := hashPassword(password)

	// bcrypt should produce different hashes for the same password (due to salt)
	if hash1 == hash2 {
		t.Error("Same password produced identical hashes (salt not working)")
	}
}

func TestCheckPasswordCorrect(t *testing.T) {
	password := "correctpassword"
	hash := hashPassword(password)

	if !checkPassword(password, hash) {
		t.Error("Valid password failed to authenticate")
	}
}

func TestCheckPasswordIncorrect(t *testing.T) {
	password := "correctpassword"
	hash := hashPassword(password)

	if checkPassword("wrongpassword", hash) {
		t.Error("Invalid password incorrectly authenticated")
	}
}

func TestCheckPasswordEmptyString(t *testing.T) {
	password := "somepassword"
	hash := hashPassword(password)

	if checkPassword("", hash) {
		t.Error("Empty password should not authenticate")
	}
}

func TestCheckPasswordCaseSensitive(t *testing.T) {
	password := "Password123"
	hash := hashPassword(password)

	// Password should be case-sensitive
	if checkPassword("password123", hash) {
		t.Error("Password check should be case-sensitive")
	}

	if checkPassword("PASSWORD123", hash) {
		t.Error("Password check should be case-sensitive")
	}
}
