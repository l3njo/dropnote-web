package controllers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/l3njo/dropnote-web/models"
	uuid "github.com/satori/go.uuid"

	"github.com/gorilla/sessions"
)

const (
	sessionCookie = "session-cookie"
	site          = "https://drop-note.herokuapp.com/"
)

var (
	base  = filepath.Join("templates", "base.html.tmpl")
	key   = []byte(os.Getenv("AES_KEY"))
	store = sessions.NewCookieStore(key)
)

// Data stores page information
type Data struct {
	Title, Site, Link, Description string
}

// Flash stores flash messages
type Flash struct {
	Message string
	Status, Custom  bool
}

// Info holds data for informational pages
type Info struct {
	Heading, Message string
}

// Page holds page data
type Page struct {
	Data
	Info
	Flashes []Flash
	models.User
	models.Note
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

func isUUID(s string) bool {
	if _, e := uuid.FromString(s); e != nil {
		return false
	}
	return true
}

// Handle deals with top-level errors
func Handle(e error) {
	if e != nil {
		log.Println(e)
	}
}
