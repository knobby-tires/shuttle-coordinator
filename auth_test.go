package main

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestLoginHandlerGET(t *testing.T) {
	initTemplates()

	req := httptest.NewRequest("GET", "/login", nil)
	w := httptest.NewRecorder()

	loginHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Jacob's Flight Tracker") {
		t.Error("Login page doesn't contain expected title")
	}
}

func TestLoginHandlerPOSTSuccess(t *testing.T) {
	initTemplates()

	form := url.Values{}
	form.Add("username", "demo")
	form.Add("password", "demo123")

	req := httptest.NewRequest("POST", "/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	loginHandler(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("Expected redirect (303), got %d", w.Code)
	}

	cookies := w.Result().Cookies()
	found := false
	for _, cookie := range cookies {
		if cookie.Name == "session_id" {
			found = true
			if cookie.HttpOnly != true {
				t.Error("Session cookie should be HttpOnly")
			}
			if cookie.Value == "" {
				t.Error("Session cookie value is empty")
			}
		}
	}
	if !found {
		t.Error("No session cookie was set")
	}
}

func TestLoginHandlerPOSTFailure(t *testing.T) {
	initTemplates()

	form := url.Values{}
	form.Add("username", "demo")
	form.Add("password", "wrongpassword")

	req := httptest.NewRequest("POST", "/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	loginHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Invalid username or password") {
		t.Error("Error message not displayed on failed login")
	}
}

func TestRequireAuthMiddleware(t *testing.T) {
	called := false
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}

	protectedHandler := requireAuth(testHandler)

	// Test without session - should redirect
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	protectedHandler(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("Expected redirect without auth, got %d", w.Code)
	}
	if called {
		t.Error("Protected handler was called without authentication")
	}

	// Test with valid session
	called = false
	sessionID := createSession("testuser")
	req = httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: sessionID,
	})
	w = httptest.NewRecorder()
	protectedHandler(w, req)

	if !called {
		t.Error("Protected handler was not called with valid session")
	}
}

func TestGetCurrentUser(t *testing.T) {
	sessionID := createSession("valet")

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: sessionID,
	})

	user := getCurrentUser(req)
	if user == nil {
		t.Fatal("getCurrentUser returned nil for valid session")
	}

	if user.Username != "valet" {
		t.Errorf("Expected username 'valet', got '%s'", user.Username)
	}

	if user.Role != "valet" {
		t.Errorf("Expected role 'valet', got '%s'", user.Role)
	}
}

func initTemplates() {
	var err error
	tmpl, err = template.New("index").Parse(htmlTemplate)
	if err != nil {
		panic(err)
	}
	tmpl, err = tmpl.New("login").Parse(loginTemplate)
	if err != nil {
		panic(err)
	}
}
