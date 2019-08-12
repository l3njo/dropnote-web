package controllers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

// MeHandler handles the "/me" route.
func MeHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionCookie)
	Handle(err)
	isAuth := checkAuth(session)
	data := info{Title: "Profile"}
	meta := filepath.Join("templates", "meta", "me.html.tmpl")
	menu := getMenu(isAuth)
	body := filepath.Join("templates", "me.html.tmpl")
	if !isAuth {
		data.Title, data.Heading = "401", "You can't access this!"
		data.Message = "You are not allowed to access this content. Try signing up or logging in."
		meta = filepath.Join("templates", "meta", "error.html.tmpl")
		body = filepath.Join("templates", "error.html.tmpl")
	} else {
		sData := session.Values["data"].(*SessionData)
		data.Name, data.Mail = sData.Name, sData.Mail
		data.Notes, err = getNotes(sData.Auth)
		Handle(err)
	}

	tmpl, err := template.New("me.html").Funcs(funcMap).ParseFiles(base, menu, meta, body)
	Handle(err)
	err = tmpl.ExecuteTemplate(w, "layout", data)
	Handle(err)
	return
}
