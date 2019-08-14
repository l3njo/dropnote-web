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
		data.Title, data.Heading = "403", "You can't access this!"
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

// NoteActionsHandler handles toggle and delete actions for notes
func NoteActionsHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionCookie)
	Handle(err)
	data := info{Title: "403"}
	if checkAuth(session) {
		uData := session.Values["data"].(*models.User)
		if actions, ok := r.URL.Query()["a"]; ok && len(actions) > 0 {
			if keys, ok := r.URL.Query()["note"]; ok && len(keys) > 0 {
				action := actions[0]
				note := &models.Note{Voucher: keys[0]}
				Handle(note.Get(uData.Auth))
				switch action {
				case "toggle":
					Handle(note.Toggle(uData.Auth))
				case "delete":
					Handle(note.Delete(uData.Auth))
				}
			}
			http.Redirect(w, r, "/me", http.StatusFound)
		}
	}

	data.Title, data.Heading = "403", "You can't access this!"
	data.Message = "You are not allowed to access this content. Try signing up or logging in."
	meta := filepath.Join("templates", "meta", "error.html.tmpl")
	body := filepath.Join("templates", "error.html.tmpl")
	tmpl, err := template.New("me.html").Funcs(funcMap).ParseFiles(base, meta, body)
	Handle(err)
	err = tmpl.ExecuteTemplate(w, "layout", data)
	Handle(err)
	return
}
