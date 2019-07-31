package app

import (
	"log"
	"net/http"
	"strings"

	c "github.com/l3njo/dropnote-web/controllers"
)

// App represents the application
type App struct {
	mux *http.ServeMux
}

// Init sets up the router
func (a *App) Init() {
	a.mux = http.NewServeMux()
	a.mux.HandleFunc("/favicon.ico", c.FaviconHandler)
	a.mux.HandleFunc("/home", c.IndexHandler)
	a.mux.HandleFunc("/", c.IndexHandler)

	fs := http.FileServer(http.Dir("static"))
	a.mux.Handle("/static/", http.StripPrefix("/static/", fs))
}

// Run starts the app
func (a *App) Run(port string) {
	log.Printf("Starting server on port %s.\n", port)
	log.Fatal("ListenAndServe: ", http.ListenAndServe(":"+port, loggingMiddleware(trailingSlashesMiddleware(a.mux))))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("RequestURI:", r.RequestURI)
		next.ServeHTTP(w, r)
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
