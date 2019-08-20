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
		displayHTTPError(w, r, http.StatusForbidden)
		return
	}
	user := session.Values["data"].(*models.User)
	data.Name, data.Mail = user.Name, user.Mail
	data.Notes, err = user.GetNotes()
	Handle(err)
	tmpl, err := template.New("me.html").Funcs(funcMap).ParseFiles(base, meta, body)
	Handle(err)
	err = tmpl.ExecuteTemplate(w, "layout", data)
	Handle(err)
	return
}

// MyActionsHandler handles the "/me/actions" route.
func MyActionsHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionCookie)
	Handle(err)
	isAuth := checkAuth(session)
	if !isAuth {
		displayHTTPError(w, r, http.StatusForbidden)
		return
	}

	if actions, ok := r.URL.Query()["a"]; ok && len(actions) > 0 {
		user := session.Values["data"].(*models.User)
		switch actions[0] {
		case "delete":
			if err := user.TryDelete(); err != nil {
				Handle(err)
				session.AddFlash(Flash{Message: "Account delete failed."})
				Handle(session.Save(r, w))
				displayHTTPError(w, r, http.StatusInternalServerError)
				return
			}
		default:
			displayHTTPError(w, r, http.StatusNotFound)
			return
		}
	}

	session.AddFlash(Flash{Message: "Account deleted.", Status: true})
	Handle(session.Save(r, w))
	http.Redirect(w, r, "/logout", http.StatusFound)
	return
}
