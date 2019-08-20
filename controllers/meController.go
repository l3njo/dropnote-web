package controllers

import (
	"fmt"
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

	if flashes := session.Flashes(); len(flashes) > 0 {
		for _, v := range flashes {
			data.Flashes = append(data.Flashes, *v.(*Flash))
		}
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
	Handle(session.Save(r, w))
	tmpl, err := template.New("me.html").Funcs(funcMap).ParseFiles(base, meta, body)
	Handle(err)
	err = tmpl.ExecuteTemplate(w, "layout", data)
	Handle(err)
	return
}

// MyUpdateHandler handles the "/me/update" route.
func MyUpdateHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionCookie)
	Handle(err)
	isAuth := checkAuth(session)
	if !isAuth {
		displayHTTPError(w, r, http.StatusForbidden)
		return
	}

	user := session.Values["data"].(*models.User)
	data := Page{
		Data: Data{
			Site: site,
			Link: r.URL.Path,
		},
		User: *user,
	}

	if flashes := session.Flashes(); len(flashes) > 0 {
		for _, v := range flashes {
			data.Flashes = append(data.Flashes, *v.(*Flash))
		}
	}

	meta := filepath.Join("templates", "meta", "auth.html.tmpl")
	body := filepath.Join("templates", "update.html.tmpl")
	context := ""
	if values, ok := r.URL.Query()["v"]; ok && len(values) > 0 {
		context = values[0]
		switch context {
		case "mail":
			data.Title, data.Description = "Change Email", "Change user email."
		case "pass":
			data.Title, data.Description = "Change Password", "Change user password."
		default:
			displayHTTPError(w, r, http.StatusNotFound)
			return
		}
	}

	if r.Method == "POST" {
		r.ParseForm()
		next := fmt.Sprintf("%s?v=%s", r.URL.Path, context)
		switch context {
		case "mail":
			mail := user.Mail
			user.Mail = r.Form["mail"][0]
			if mail == user.Mail {
				session.AddFlash(Flash{Message: "choose a different email", Status: warning, Custom: true})
			} else if err := user.TryMailUpdate(); err != nil {
				Handle(err)
				session.AddFlash(Flash{Message: err.Error(), Status: warning, Custom: true})
			} else {
				next = "/me"
				session.Values["data"] = user
				session.AddFlash(Flash{Message: "Update successful.", Status: success})
			}
		case "pass":
			user.Current, user.Updated, user.Confirm = r.Form["current"][0], r.Form["updated"][0], r.Form["confirm"][0]
			if err := user.ValidatePassUpdate(); err != nil {
				Handle(err)
				session.AddFlash(Flash{Message: err.Error(), Status: warning, Custom: true})
			} else if err := user.TryPassUpdate(); err != nil {
				Handle(err)
				session.AddFlash(Flash{Message: err.Error(), Status: warning, Custom: true})
			} else {
				next = "/me"
				session.AddFlash(Flash{Message: "Update successful.", Status: success})
			}
		default:
			displayHTTPError(w, r, http.StatusNotFound)
			return
		}

		Handle(session.Save(r, w))
		http.Redirect(w, r, next, http.StatusFound)
		return
	}

	Handle(session.Save(r, w))
	tmpl, err := template.ParseFiles(base, meta, body)
	Handle(err)
	tmpl.ExecuteTemplate(w, "layout", data)
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
				session.AddFlash(Flash{Message: "Account delete failed.", Status: warning})
				Handle(session.Save(r, w))
				displayHTTPError(w, r, http.StatusInternalServerError)
				return
			}
		default:
			displayHTTPError(w, r, http.StatusNotFound)
			return
		}
	}

	session.AddFlash(Flash{Message: "Account deleted.", Status: success})
	Handle(session.Save(r, w))
	http.Redirect(w, r, "/logout", http.StatusFound)
	return
}
