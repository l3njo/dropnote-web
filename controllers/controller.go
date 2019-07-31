package controllers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func handle(e error) {
	if e != nil {
		log.Println(e)
	}
}

// FaviconHandler handles the "/favicon.ico" route.
func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/img/favicon.ico")
}

// IndexHandler handles the "/" and "/home" routes.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/home" {
		NotFoundHandler(w, r)
		return
	}

	base := filepath.Join("templates", "base.html.tmpl")
	nav := filepath.Join("templates", "navbar.html.tmpl")
	style := filepath.Join("templates", "styles", "home.html.tmpl")
	content := filepath.Join("templates", "home.html.tmpl")

	tmpl, err := template.ParseFiles(base, nav, style, content)
	handle(err)
	tmpl.ExecuteTemplate(w, "layout", "Home")
	return
}

// NotFoundHandler handles HTTP 404 errors.
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	base := filepath.Join("templates", "base.html.tmpl")
	nav := filepath.Join("templates", "navbar.html.tmpl")
	style := filepath.Join("templates", "styles", "404.html.tmpl")
	content := filepath.Join("templates", "404.html.tmpl")

	tmpl, err := template.ParseFiles(base, nav, style, content)
	handle(err)
	tmpl.ExecuteTemplate(w, "layout", "404")
}
