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

type SessionData struct {
	ID, User, Name, Mail, Auth string
}
