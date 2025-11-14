package main

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
	"os"

	"golang.org/x/crypto/bcrypt"
)

// User represents a user account
type User struct {
	Username     string
	PasswordHash string
	Role         string // "valet", "desk", or "demo"
}

// In-memory session storage (in production, use Redis or similar)
var (
	sessions = make(map[string]string) // sessionID -> username
	sessLock sync.RWMutex
)

// Hardcoded users with bcrypt-hashed passwords
var users = map[string]User{
	"valet": {
		Username:     "valet",
		PasswordHash: hashPassword(os.Getenv("VALET_PASSWORD")),
		Role:         "valet",
	},
	"desk": {
		Username:     "desk",
		PasswordHash: hashPassword(os.Getenv("DESK_PASSWORD")),
		Role:         "desk",
	},
	"demo": {
		Username:     "demo",
		PasswordHash: hashPassword(os.Getenv("DEMO_PASSWORD")),
		Role:         "demo",
	},
}

// hashPassword creates a bcrypt hash of a password
func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}

// checkPassword verifies a password against a hash
func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// generateSessionID creates a random session ID
func generateSessionID() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// createSession creates a new session for a user
func createSession(username string) string {
	sessLock.Lock()
	defer sessLock.Unlock()

	sessionID := generateSessionID()
	sessions[sessionID] = username
	return sessionID
}

// getSession returns the username for a session ID, or empty string if invalid
func getSession(sessionID string) string {
	sessLock.RLock()
	defer sessLock.RUnlock()

	return sessions[sessionID]
}

// deleteSession removes a session
func deleteSession(sessionID string) {
	sessLock.Lock()
	defer sessLock.Unlock()

	delete(sessions, sessionID)
}

// getCurrentUser returns the User object for the current session
func getCurrentUser(r *http.Request) *User {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil
	}

	username := getSession(cookie.Value)
	if username == "" {
		return nil
	}

	user, exists := users[username]
	if !exists {
		return nil
	}

	return &user
}

// requireAuth is middleware that checks if user is authenticated
func requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		username := getSession(cookie.Value)
		if username == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next(w, r)
	}
}

// loginHandler displays login page and handles login attempts
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Check if already logged in
		cookie, err := r.Cookie("session_id")
		if err == nil && getSession(cookie.Value) != "" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		tmpl.ExecuteTemplate(w, "login", nil)
		return
	}

	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		user, exists := users[username]

		if !exists || !checkPassword(password, user.PasswordHash) {
			tmpl.ExecuteTemplate(w, "login", map[string]string{
				"Error": "Invalid username or password",
			})
			return
		}

		sessionID := createSession(username)

		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			MaxAge:   86400 * 7, // 7 days
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// logoutHandler destroys the session and logs out the user
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		deleteSession(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
