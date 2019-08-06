package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

const(
	site = "https://drop-note.herokuapp.com/"
	api = "https://dropnote-api.herokuapp.com/api/"
)

var(
	base = filepath.Join("templates", "base.html.tmpl")
)

type info struct {
	Title, Heading, Message string
}

func handle(e error) {
	if e != nil {
		log.Println(e)
	}
}

// FaviconHandler handles the "/favicon.ico" route.
func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/img/favicon.ico")
}

// IndexHandler handles the "/", "/home" and all undefined routes.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	data := info{Title: "Home"}
	nav := filepath.Join("templates", "navbar.html.tmpl")
	style := filepath.Join("templates", "styles", "home.html.tmpl")
	content := filepath.Join("templates", "home.html.tmpl")
	if r.URL.Path != "/" && r.URL.Path != "/home" {
		data.Title = "404"
		style = filepath.Join("templates", "styles", "404.html.tmpl")
		content = filepath.Join("templates", "404.html.tmpl")
	}

	tmpl, err := template.ParseFiles(base, nav, style, content)
	handle(err)
	tmpl.ExecuteTemplate(w, "layout", data)
	return
}

// DropNoteHandler handles the "/dropnote" route.
func DropNoteHandler(w http.ResponseWriter, r *http.Request) {
	data := info{Title: "Drop Note"}
	nav := filepath.Join("templates", "navbar.html.tmpl")
	style := filepath.Join("templates", "styles", "drop.html.tmpl")
	content := filepath.Join("templates", "dropnote.html.tmpl")
	if r.Method == "POST" {
		r.ParseForm()
		subject := r.Form["subject"][0]
		payload := note{
			subject: subject, 
			content: r.Form["content"][0],
		} 
		url := fmt.Sprintf("%snote/new", api)
		data.Heading, data.Message = "Error!", "Spmething has gone horribly wrong"
		if ok, voucher := postNote(url, payload); ok {
			strOnSuccess := "Your note (%s) has been stored.\nYour code is %s.\nHere's a direct link: %sdropcode?voucher=%s."
			data.Heading, data.Message = "Success!", fmt.Sprintf(strOnSuccess, subject, voucher, site, voucher)
		} 
		style = filepath.Join("templates", "styles", "info.html.tmpl")
		content = filepath.Join("templates", "info.html.tmpl")
	}

	tmpl, err := template.ParseFiles(base, nav, style, content)
	handle(err)
	tmpl.ExecuteTemplate(w, "layout", data)
	return
}

// DropCodeHandler handles the "/dropcode" route.
func DropCodeHandler(w http.ResponseWriter, r *http.Request) {
	var voucher, content, style string
	data := info{Title: "Drop Code"}
	nav := filepath.Join("templates", "navbar.html.tmpl")
	style = filepath.Join("templates", "styles", "drop.html.tmpl")
	content = filepath.Join("templates", "dropcode.html.tmpl")
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
		if validate(voucher){
			url := fmt.Sprintf("%snote/%s", api, voucher)
			data.Message = "Something has gone horribly wrong."
			if noteData, err := getNote(url); err == nil {
				data.Heading, data.Message = noteData.subject, noteData.content
			}
		}
		style = filepath.Join("templates", "styles", "info.html.tmpl")
		content = filepath.Join("templates", "info.html.tmpl")
	} 
	
	tmpl, err := template.ParseFiles(base, nav, style, content)
	handle(err)
	tmpl.ExecuteTemplate(w, "layout", data)
	return
}