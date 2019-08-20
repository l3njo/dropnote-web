package app

import (
	"log"
	"net/http"
	"strings"
	"time"

	c "github.com/l3njo/dropnote-web/controllers"
)

// Application represents the application
type Application struct {
	mux *http.ServeMux
}

// Init sets up the router
func (a *Application) Init() {
	a.mux = http.NewServeMux()
	a.mux.HandleFunc("/home", c.IndexHandler)
	a.mux.HandleFunc("/dropnote", c.DropNoteHandler)
	a.mux.HandleFunc("/dropcode", c.DropCodeHandler)

	a.mux.HandleFunc("/signup", c.SignupHandler)
	a.mux.HandleFunc("/login", c.LoginHandler)
	a.mux.HandleFunc("/logout", c.LogoutHandler)
	a.mux.HandleFunc("/reset", c.ResetHandler)

	a.mux.HandleFunc("/me/notes/update/", c.NoteUpdateHandler)
	a.mux.HandleFunc("/me/notes/action", c.NoteActionsHandler)
	a.mux.HandleFunc("/me/notes/", c.MyNotesHandler)
	a.mux.HandleFunc("/me/action", c.MyActionsHandler)
	a.mux.HandleFunc("/me", c.MeHandler)

	a.mux.HandleFunc("/temp/qr/notes/", c.QRHandler)
	a.mux.HandleFunc("/", c.IndexHandler)

	fs := http.FileServer(http.Dir("static"))
	a.mux.Handle("/static/", http.StripPrefix("/static/", fs))
}

// Run starts the app
func (a *Application) Run(port string) {
	log.Printf("Starting server on port %s.\n", port)
	log.Fatal("ListenAndServe: ", http.ListenAndServe(":"+port, loggingMiddleware(trailingSlashesMiddleware(a.mux))))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%-6s%-57s\t%s", r.Method, r.RequestURI, time.Since(start))
	})
}

func trailingSlashesMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}
		next.ServeHTTP(w, r)
	})
}
