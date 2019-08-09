package controllers

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/sessions"
)

const (
	site = "https://drop-note.herokuapp.com/"
	api  = "https://dropnote-api.herokuapp.com/api/"
)

var (
	base = filepath.Join("templates", "base.html.tmpl")
)

type info struct {
	Title, Heading, Message, User string
}

func checkAuth(session *sessions.Session) bool {
	if auth, ok := session.Values["isAuth"].(bool); !ok || !auth {
		return false
	}
	return true
}

func getMenu(isAuth bool) string {
	if isAuth {
		return filepath.Join("templates", "menu", "private_nav.html.tmpl")
	}
	return filepath.Join("templates", "menu", "public_nav.html.tmpl")
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
