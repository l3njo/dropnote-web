package controllers

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/l3njo/dropnote-web/models"
)

// MeHandler handles the "/me" route.
func MeHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionCookie)
	Handle(err)
	isAuth := checkAuth(session)
	data := info{Title: "Profile"}
	meta := filepath.Join("templates", "meta", "me.html.tmpl")
	body := filepath.Join("templates", "me.html.tmpl")
	if !isAuth {
		data.Title, data.Heading = "401", "You can't access this!"
		data.Message = "You are not allowed to access this content. Try signing up or logging in."
		meta = filepath.Join("templates", "meta", "error.html.tmpl")
		body = filepath.Join("templates", "error.html.tmpl")
	} else {
		user := session.Values["data"].(*models.User)
		data.Name, data.Mail = user.Name, user.Mail
		data.Notes, err = user.GetNotes()
		Handle(err)
	}

	tmpl, err := template.New("me.html").Funcs(funcMap).ParseFiles(base, meta, body)
	Handle(err)
	err = tmpl.ExecuteTemplate(w, "layout", data)
	Handle(err)
	return
}
