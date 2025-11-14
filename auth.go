package main

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"os"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

// User represents an authenticated account with role-based access
type User struct {
	Username     string
	PasswordHash string
	Role         string // "valet", "desk", or "demo"
}

// Session storage with mutex for concurrent access
var (
	sessions = make(map[string]string) // sessionID -> username
	sessLock sync.RWMutex
)

// Hardcoded users - in production these would be in a database
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

// hashPassword creates a bcrypt hash from plain text password
func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}

// checkPassword verifies a password against a bcrypt hash
func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// generateSessionID creates a cryptographically secure random session ID
func generateSessionID() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// createSession generates and stores a new session for a user
func createSession(username string) string {
	sessLock.Lock()
	defer sessLock.Unlock()

	sessionID := generateSessionID()
	sessions[sessionID] = username
	return sessionID
}

// getSession retrieves username from session ID
func getSession(sessionID string) string {
	sessLock.RLock()
	defer sessLock.RUnlock()

	return sessions[sessionID]
}

// deleteSession removes a session (logout)
func deleteSession(sessionID string) {
	sessLock.Lock()
	defer sessLock.Unlock()

	delete(sessions, sessionID)
}

// getCurrentUser extracts User from request session cookie
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

// requireAuth is middleware that protects routes requiring authentication
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

// loginHandler displays login form and processes login attempts
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Redirect if already logged in
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

		// Create session and set secure cookie
		sessionID := createSession(username)

		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			MaxAge:   86400 * 7, // 7 days
			HttpOnly: true,      // XSS protection
			SameSite: http.SameSiteStrictMode, // CSRF protection
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// logoutHandler destroys session and clears cookie
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
