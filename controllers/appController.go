package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/l3njo/dropnote-web/models"
	uuid "github.com/satori/go.uuid"
)

var funcMap = template.FuncMap{
	"inc": func(i int) int {
		return i + 1
	},
}

// IndexHandler handles the "/", "/home", "/favicon.ico" and all undefined routes.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		http.ServeFile(w, r, "static/img/favicon.ico")
		return
	}

	session, err := store.Get(r, sessionCookie)
	Handle(err)
	isAuth := checkAuth(session)
	data := info{Title: "Home"}
	meta := filepath.Join("templates", "meta", "home.html.tmpl")
	body := filepath.Join("templates", "home.html.tmpl")
	if isAuth {
		uData := session.Values["data"].(*models.User)
		data.Name, data.Mail = uData.Name, uData.Mail
	}

	if r.URL.Path != "/" && r.URL.Path != "/home" {
		data.Title, data.Heading = "404", "We are sorry, Page not found!"
		data.Message = "The page you are looking for might have been removed had its name changed or is temporarily unavailable."
		meta = filepath.Join("templates", "meta", "error.html.tmpl")
		body = filepath.Join("templates", "error.html.tmpl")
	}

	tmpl, err := template.ParseFiles(base, meta, body)
	Handle(err)
	tmpl.ExecuteTemplate(w, "layout", data)
	return
}

// DropNoteHandler handles the "/dropnote" route.
func DropNoteHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionCookie)
	Handle(err)
	isAuth := checkAuth(session)
	data := info{Title: "Drop Note"}
	meta := filepath.Join("templates", "meta", "drop.html.tmpl")
	body := filepath.Join("templates", "dropnote.html.tmpl")
	if isAuth {
		uData := session.Values["data"].(*models.User)
		data.Name, data.Mail = uData.Name, uData.Mail
	}

	if r.Method == "POST" {
		r.ParseForm()
		note := &models.Note{
			Subject: r.Form["subject"][0],
			Content: r.Form["content"][0],
		}

		auth := ""
		if isAuth && len(r.Form["shouldLink"]) > 0 {
			if shouldLink := r.Form["shouldLink"][0]; shouldLink == "on" {
				sessionData := session.Values["data"].(*models.User)
				auth = sessionData.Auth
			}
		}

		if err := note.Post(auth); err != nil {
			data.Heading, data.Message = "Error!", err.Error()
		} else {
			voucher := note.Voucher
			strOnSuccess := "Your note (%s) has been stored. Your code is %s. Here's a direct link: %sdropcode?voucher=%s."
			data.Heading, data.Message = "Success!", fmt.Sprintf(strOnSuccess, note.Subject, voucher, site, voucher)
		}
		meta = filepath.Join("templates", "meta", "info.html.tmpl")
		body = filepath.Join("templates", "info.html.tmpl")
	}

	tmpl, err := template.ParseFiles(base, meta, body)
	Handle(err)
	tmpl.ExecuteTemplate(w, "layout", data)
	return
}

// DropCodeHandler handles the "/dropcode" route.
func DropCodeHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionCookie)
	Handle(err)
	isAuth := checkAuth(session)
	uData := &models.User{}
	data := info{Title: "Drop Code"}
	meta := filepath.Join("templates", "meta", "drop.html.tmpl")
	body := filepath.Join("templates", "dropcode.html.tmpl")
	if isAuth {
		uData = session.Values["data"].(*models.User)
		data.Name, data.Mail = uData.Name, uData.Mail
	}

	if keys, ok := r.URL.Query()["voucher"]; ok && len(keys) > 0 {
		note := &models.Note{Voucher: keys[0]}
		if err := note.ValidateGet(); err != nil {
			data.Heading, data.Message = "Error!", err.Error()
		} else {
			if err := note.Get(uData.Auth); err != nil {
				data.Message = err.Error()
			} else if *note == (models.Note{}) {
				data.Message = "That note does not exist"
			} else if uuid.Equal(uuid.FromStringOrNil(note.Voucher), uuid.Nil) {
				data.Message = "That note does not exist"
			} else {
				data.Heading, data.Message = note.Subject, note.Content
			}
		}
		meta = filepath.Join("templates", "meta", "info.html.tmpl")
		body = filepath.Join("templates", "info.html.tmpl")
	}

	tmpl, err := template.ParseFiles(base, meta, body)
	Handle(err)
	tmpl.ExecuteTemplate(w, "layout", data)
	return
}
