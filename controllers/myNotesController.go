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
		displayHTTPError(w, r, http.StatusForbidden)
		return
	}
	user := session.Values["data"].(*models.User)
	code := strings.TrimPrefix(r.URL.Path, "/me/notes/")
	note := &models.Note{Voucher: code}
	data.Name, data.Mail, data.Description = user.Name, user.Mail, fmt.Sprintf(data.Description, code)
	Handle(note.Get(user.Auth))
	Handle(note.ParseDate())
	data.Note = *note
	session.AddFlash(Flash{Message: "Note retrieved.", Status: success})
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
		displayHTTPError(w, r, http.StatusForbidden)
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
					session.AddFlash(Flash{Message: "No fields have been changed", Status: warning})
					Handle(session.Save(r, w))
					http.Redirect(w, r, "/me", http.StatusFound)
					return
				}
				note.Creator = uData.User
			}
		}

		note.Subject, note.Content = subject, content
		if err := note.Update(uData.Auth); err != nil {
			Handle(err)
			session.AddFlash(Flash{Message: "Update failed.", Status: warning})
			Handle(session.Save(r, w))
			displayHTTPError(w, r, http.StatusInternalServerError)
			return
		}
		session.AddFlash(Flash{Message: "Note saved.", Status: success})
		Handle(session.Save(r, w))
		http.Redirect(w, r, "/me", http.StatusFound)
		return
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

	displayHTTPError(w, r, http.StatusForbidden)
	return
}
