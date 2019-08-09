package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// IndexHandler handles the "/", "/home", "/favicon.ico" and all undefined routes.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionCookie)
	isAuth := checkAuth(session)
	log.Println(isAuth) // DEBUG
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
	handle(err)
	tmpl.ExecuteTemplate(w, "layout", data)
	return
}

// DropNoteHandler handles the "/dropnote" route.
func DropNoteHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionCookie)
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
		payload := note{
			subject: subject,
			content: r.Form["content"][0],
		}
		url := fmt.Sprintf("%snote/new", api)
		data.Heading, data.Message = "Error!", "Something has gone horribly wrong"
		if ok, voucher := postNote(url, payload); ok {
			strOnSuccess := "Your note (%s) has been stored.\nYour code is %s.\nHere's a direct link: %sdropcode?voucher=%s."
			data.Heading, data.Message = "Success!", fmt.Sprintf(strOnSuccess, subject, voucher, site, voucher)
		}
		meta = filepath.Join("templates", "meta", "info.html.tmpl")
		body = filepath.Join("templates", "info.html.tmpl")
	}

	tmpl, err := template.ParseFiles(base, menu, meta, body)
	handle(err)
	tmpl.ExecuteTemplate(w, "layout", data)
	return
}

// DropCodeHandler handles the "/dropcode" route.
func DropCodeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionCookie)
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

	keys, hasGet := r.URL.Query()["voucher"]
	hasGet = hasGet && len(keys) > 0

	if r.Method == "POST" || hasGet {
		if hasGet {
			voucher = keys[0]
		} else {
			r.ParseForm()
			voucher = r.Form["voucher"][0]
		}

		data.Heading, data.Message = "Error!", "Your voucher is invalid."
		if validate(voucher) {
			url := fmt.Sprintf("%snote/%s", api, voucher)
			data.Message = "Something has gone horribly wrong."
			if noteData, err := getNote(url); err == nil {
				data.Heading, data.Message = noteData.subject, noteData.content
			}
		}
		meta = filepath.Join("templates", "meta", "info.html.tmpl")
		body = filepath.Join("templates", "info.html.tmpl")
	}

	tmpl, err := template.ParseFiles(base, menu, meta, body)
	handle(err)
	tmpl.ExecuteTemplate(w, "layout", data)
	return
}
