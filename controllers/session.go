package controllers

import (
	"os"

	"github.com/gorilla/sessions"
)

const (
	sessionCookie = "session-cookie"
)

var (
	key   = []byte(os.Getenv("AES_KEY"))
	store = sessions.NewCookieStore(key)
)

// SessionData represents the data stored by session cookie
type SessionData struct {
	User, Name, Mail, Auth string
}

func (s *SessionData) buildSession(result map[string]interface{}) {
	resultData := result["data"].(map[string]interface{})
	s.User = resultData["ID"].(string)
	s.Name = resultData["name"].(string)
	s.Mail = resultData["mail"].(string)
	s.Auth = resultData["auth"].(string)
	return
}
