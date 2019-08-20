package controllers

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/l3njo/dropnote-web/models"
	uuid "github.com/satori/go.uuid"

	"github.com/gorilla/sessions"
)

type flashStatus string

const (
	info    flashStatus = "info"
	warning flashStatus = "warning"
	danger  flashStatus = "danger"
	success flashStatus = "success"
)

// Data stores page information
type Data struct {
	Title, Site, Link, Description string
}

// Flash stores flash messages
type Flash struct {
	Message string
	Status  flashStatus
	Custom  bool
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

const (
	sessionCookie = "session-cookie"
	site          = "https://drop-note.herokuapp.com"
)

var (
	base       = filepath.Join("templates", "base.html.tmpl")
	key        = []byte(os.Getenv("AES_KEY"))
	store      = sessions.NewCookieStore(key)
	httpErrors = map[int]Info{
		http.StatusNotFound: Info{
			Heading: "Sorry, we can't find that.",
			Message: "The page you are looking for might have been removed or is temporarily unavailable.",
		},
		http.StatusForbidden: Info{
			Heading: "Hey! You shouldn't be here!",
			Message: "You are not allowed to access this content. Try signing up or logging in.",
		},
		http.StatusInternalServerError: Info{
			Heading: "Oops! Something has gone horribly wrong.",
			Message: "A connection to the server couldn't be established. Please try again later.",
		},
	}
)

func checkAuth(s *sessions.Session) bool {
	if auth, ok := s.Values["isAuth"].(bool); !ok || !auth {
		return false
	}
	return true
}

func displayHTTPError(w http.ResponseWriter, r *http.Request, e int) {
	data := Page{
		Data: Data{
			Title: strconv.Itoa(e),
		},
	}

	info, ok := httpErrors[e]
	if !ok {
		info = httpErrors[http.StatusInternalServerError]
		data.Title = strconv.Itoa(e)
	}

	data.Info = info
	meta := filepath.Join("templates", "meta", "error.html.tmpl")
	body := filepath.Join("templates", "error.html.tmpl")
	tmpl, err := template.ParseFiles(base, meta, body)
	Handle(err)
	tmpl.ExecuteTemplate(w, "layout", data)
	return
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
