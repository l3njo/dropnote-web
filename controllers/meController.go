package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/l3njo/dropnote-web/models"
	uuid "github.com/satori/go.uuid"
)

// MeHandler handles the "/me" route.
func MeHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionCookie)
	Handle(err)
	isAuth := checkAuth(session)
	data := Page{
		Data: Data{
			Title:       "Profile",
			Site:        site,
			Link:        r.URL.Path,
			Description: "My DropNote Profile",
		},
	}
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

// MyNotesHandler handles displaying details for one note
func MyNotesHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionCookie)
	Handle(err)
	isAuth := checkAuth(session)
	data := Page{
		Data: Data{
			Title:       "Note Details",
			Site:        site,
			Link:        r.URL.Path,
			Description: "View note details.",
		},
	}
	meta := filepath.Join("templates", "meta", "note.html.tmpl")
	body := filepath.Join("templates", "note.html.tmpl")

	if flashes := session.Flashes(); len(flashes) > 0 {
		for _, v := range flashes {
			data.Flashes = append(data.Flashes, *v.(*Flash))
		}
	}

	if !isAuth {
		data.Title, data.Heading = "403", "You can't access this!"
		data.Message = "You are not allowed to access this content. Try signing up or logging in."
		meta = filepath.Join("templates", "meta", "error.html.tmpl")
		body = filepath.Join("templates", "error.html.tmpl")
	} else {
		user := session.Values["data"].(*models.User)
		data.Name, data.Mail = user.Name, user.Mail
		code := strings.TrimPrefix(r.URL.Path, "/me/notes/")
		data.Description = fmt.Sprintf(data.Description, code)
		note := &models.Note{Voucher: code}
		Handle(note.Get(user.Auth))
		Handle(note.ParseDate())
		data.Note = *note
		Handle(err)
		session.AddFlash(Flash{Message: "Note retrieved.", Status: true})
	}

	Handle(session.Save(r, w))
	tmpl, err := template.ParseFiles(base, meta, body)
	Handle(err)
	tmpl.ExecuteTemplate(w, "layout", data)
	return
}

// NoteUpdateHandler takes care of updating notes
func NoteUpdateHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionCookie)
	Handle(err)
	isAuth := checkAuth(session)
	data := Page{
		Data: Data{
			Title:       "Edit Note",
			Site:        site,
			Link:        r.URL.Path,
			Description: "Edit a Note",
		},
	}
	meta := filepath.Join("templates", "meta", "drop.html.tmpl")
	body := filepath.Join("templates", "dropnote.html.tmpl")
	if !isAuth {
		data.Title, data.Heading = "403", "You can't access this!"
		data.Message = "You are not allowed to access this content. Try signing up or logging in."
		meta = filepath.Join("templates", "meta", "error.html.tmpl")
		body = filepath.Join("templates", "error.html.tmpl")
		tmpl, err := template.ParseFiles(base, meta, body)
		Handle(err)
		tmpl.ExecuteTemplate(w, "layout", data)
		return
	}

	uData := session.Values["data"].(*models.User)
	data.Name, data.Mail = uData.Name, uData.Mail
	if flashes := session.Flashes(); len(flashes) > 0 {
		for _, v := range flashes {
			data.Flashes = append(data.Flashes, *v.(*Flash))
		}
	}

	voucher := strings.TrimPrefix(r.URL.Path, "/me/notes/update/")
	note := models.Note{Voucher: voucher}
	Handle(note.Get(uData.Auth))
	data.Note = note

	if r.Method == "POST" {
		r.ParseForm()
		subject, content := r.Form["subject"][0], r.Form["content"][0]
		note.Creator = uuid.Nil.String()
		if isAuth && len(r.Form["shouldLink"]) > 0 {
			if shouldLink := r.Form["shouldLink"][0]; shouldLink == "on" {
				if subject == note.Subject {
					subject = ""
				}
				if content == note.Content {
					content = ""
				}

				if unchanged := subject == "" && content == ""; unchanged {
					session.AddFlash(Flash{Message: "No fields have been changed"})
					Handle(session.Save(r, w))
					http.Redirect(w, r, "/me", http.StatusFound)
					return
				}
				note.Creator = uData.User
			}
		}

		note.Subject, note.Content = subject, content
		if err := note.Update(uData.Auth); err != nil {
			session.AddFlash(Flash{Message: "Update failed."})
			data.Heading, data.Message = "Error!", err.Error()
			meta = filepath.Join("templates", "meta", "info.html.tmpl")
			body = filepath.Join("templates", "info.html.tmpl")
		} else {
			session.AddFlash(Flash{Message: "Note saved.", Status: true})
			Handle(session.Save(r, w))
			http.Redirect(w, r, "/me", http.StatusFound)
			return
		}
	}

	Handle(session.Save(r, w))
	tmpl, err := template.ParseFiles(base, meta, body)
	Handle(err)
	tmpl.ExecuteTemplate(w, "layout", data)
	return
}

// NoteActionsHandler handles toggle and delete actions for notes
func NoteActionsHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionCookie)
	Handle(err)
	data := Page{
		Data: Data{
			Title:       "Note Actions",
			Site:        site,
			Link:        r.URL.Path,
			Description: "Make changes to a note.",
		},
	}
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
