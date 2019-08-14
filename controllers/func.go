package controllers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/l3njo/dropnote-web/models"

	"github.com/gorilla/sessions"
)

const (
	sessionCookie = "session-cookie"
	site          = "https://drop-note.herokuapp.com/"
	api           = "https://dropnote-api.herokuapp.com/api/"
)

var (
	base  = filepath.Join("templates", "base.html.tmpl")
	key   = []byte(os.Getenv("AES_KEY"))
	store = sessions.NewCookieStore(key)
)

type info struct {
	Title, Heading, Message string
	models.User
	Notes []models.Note
}

func checkAuth(s *sessions.Session) bool {
	if auth, ok := s.Values["isAuth"].(bool); !ok || !auth {
		return false
	}
	return true
}

func getNext(r *http.Request) string {
	nextList, ok := r.URL.Query()["next"]
	if !ok || len(nextList[0]) < 1 {
		return "/"
	}
	next := nextList[0]
	return next
}

// Handle deals with top-level errors
func Handle(e error) {
	if e != nil {
		log.Println(e)
	}
}
