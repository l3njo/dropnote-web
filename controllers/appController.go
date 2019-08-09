package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

// IndexHandler handles the "/", "/home", "/favicon.ico" and all undefined routes.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionCookie)
	Handle(err)
	isAuth := checkAuth(session)
	if r.URL.Path == "/favicon.ico" {
		http.ServeFile(w, r, "static/img/favicon.ico")
		return
	}

	data := info{Title: "Home"}
	meta := filepath.Join("templates", "meta", "home.html.tmpl")
	menu := getMenu(isAuth)
	body := filepath.Join("templates", "home.html.tmpl")
	if isAuth {
		sData := session.Values["data"].(*SessionData)
		data.User = sData.Name
	}

	if r.URL.Path != "/" && r.URL.Path != "/home" {
		data.Title = "404"
		meta = filepath.Join("templates", "meta", "404.html.tmpl")
		body = filepath.Join("templates", "404.html.tmpl")
	}

	tmpl, err := template.ParseFiles(base, menu, meta, body)
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
	menu := getMenu(isAuth)
	body := filepath.Join("templates", "dropnote.html.tmpl")
	if isAuth {
		sData := session.Values["data"].(*SessionData)
		data.User = sData.Name
	}

	if r.Method == "POST" {
		r.ParseForm()
		subject := r.Form["subject"][0]
		noteData := &note{
			Subject: subject,
			Content: r.Form["content"][0],
		}

		if isAuth && len(r.Form["shouldLink"]) > 0 {
			if shouldLink := r.Form["shouldLink"][0]; shouldLink == "on" {
				sessionData := session.Values["data"].(*SessionData)
				noteData.auth = sessionData.Auth
			}
		}

		data.Heading, data.Message = "Error!", "Something has gone horribly wrong"
		if ok, voucher := noteData.postNote(); ok {
			strOnSuccess := "Your note (%s) has been stored. Your code is %s. Here's a direct link: %sdropcode?voucher=%s."
			data.Heading, data.Message = "Success!", fmt.Sprintf(strOnSuccess, subject, voucher, site, voucher)
		}
		meta = filepath.Join("templates", "meta", "info.html.tmpl")
		body = filepath.Join("templates", "info.html.tmpl")
	}

	tmpl, err := template.ParseFiles(base, menu, meta, body)
	Handle(err)
	tmpl.ExecuteTemplate(w, "layout", data)
	return
}

// DropCodeHandler handles the "/dropcode" route.
func DropCodeHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionCookie)
	Handle(err)
	isAuth := checkAuth(session)
	var voucher string
	data := info{Title: "Drop Code"}
	meta := filepath.Join("templates", "meta", "drop.html.tmpl")
	menu := getMenu(isAuth)
	body := filepath.Join("templates", "dropcode.html.tmpl")
	if isAuth {
		sData := session.Values["data"].(*SessionData)
		data.User = sData.Name
	}

	if keys, ok := r.URL.Query()["voucher"]; ok && len(keys) > 0 {
		voucher = keys[0]
		data.Heading, data.Message = "Error!", "Your voucher is invalid."
		if validateCode(voucher) {
			noteData := &note{}
			data.Message = "Something has gone horribly wrong."
			if err := noteData.getNote(voucher); err == nil {
				data.Heading, data.Message = noteData.Subject, noteData.Content
			}
		}
		meta = filepath.Join("templates", "meta", "info.html.tmpl")
		body = filepath.Join("templates", "info.html.tmpl")
	}

	tmpl, err := template.ParseFiles(base, menu, meta, body)
	Handle(err)
	tmpl.ExecuteTemplate(w, "layout", data)
	return
}
