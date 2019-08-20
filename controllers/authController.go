package controllers

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/l3njo/dropnote-web/models"
)

// SignupHandler handles the "/signup" route.
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionCookie)
	Handle(err)
	data := Page{
		Data: Data{
			Title:       "Sign Up",
			Site:        site,
			Link:        r.URL.Path,
			Description: "Sign Up for DropNote",
		},
	}
	meta := filepath.Join("templates", "meta", "auth.html.tmpl")
	body := filepath.Join("templates", "signup.html.tmpl")
	if flashes := session.Flashes(); len(flashes) > 0 {
		for _, v := range flashes {
			data.Flashes = append(data.Flashes, *v.(*Flash))
		}
		Handle(session.Save(r, w))
	}

	if r.Method == "POST" {
		r.ParseForm()
		user := &models.User{
			Name:    r.Form["name"][0],
			Mail:    r.Form["mail"][0],
			Pass:    r.Form["pass"][0],
			Confirm: r.Form["confirm"][0],
		}

		next := "/signup"
		if err := user.ValidateSignup(); err != nil {
			session.AddFlash(Flash{Message: err.Error(), Status: warning, Custom: true})
		} else if err := user.TrySignup(); err != nil {
			session.AddFlash(Flash{Message: err.Error(), Status: warning, Custom: true})
		} else {
			session.AddFlash(Flash{Message: "Signup successful.", Status: success})
			session.Values["isAuth"] = true
			session.Values["data"] = user
			next = getNext(r)
		}
		Handle(session.Save(r, w))
		http.Redirect(w, r, next, http.StatusFound)
		return
	}

	tmpl, err := template.ParseFiles(base, meta, body)
	Handle(err)
	tmpl.ExecuteTemplate(w, "layout", data)
	return
}

// LoginHandler handles the "/login" route.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionCookie)
	Handle(err)
	data := Page{
		Data: Data{
			Title:       "Log In",
			Site:        site,
			Link:        r.URL.Path,
			Description: "Log In to DropNote",
		},
	}
	meta := filepath.Join("templates", "meta", "auth.html.tmpl")
	body := filepath.Join("templates", "login.html.tmpl")
	if flashes := session.Flashes(); len(flashes) > 0 {
		for _, v := range flashes {
			data.Flashes = append(data.Flashes, *v.(*Flash))
		}
		Handle(session.Save(r, w))
	}

	if r.Method == "POST" {
		r.ParseForm()
		user := &models.User{
			Mail: r.Form["mail"][0],
			Pass: r.Form["pass"][0],
		}

		next := "/login"
		if err := user.TryLogin(); err != nil {
			session.AddFlash(Flash{Message: err.Error(), Status: warning, Custom: true})
		} else {
			session.AddFlash(Flash{Message: "Login successful.", Status: success})
			session.Values["isAuth"] = true
			session.Values["data"] = user
			next = getNext(r)
		}
		Handle(session.Save(r, w))
		http.Redirect(w, r, next, http.StatusFound)
		return
	}

	tmpl, err := template.ParseFiles(base, meta, body)
	Handle(err)
	tmpl.ExecuteTemplate(w, "layout", data)
	return
}

// LogoutHandler handles the "/logout" route.
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionCookie)
	Handle(err)
	session.Values["isAuth"] = false
	delete(session.Values, "data")
	Handle(session.Save(r, w))
	http.Redirect(w, r, getNext(r), http.StatusFound)
	return
}

// ResetHandler handles the "/reset" route.
func ResetHandler(w http.ResponseWriter, r *http.Request) {
	data := Page{
		Data: Data{
			Title:       "Reset Password",
			Site:        site,
			Link:        r.URL.Path,
			Description: "DropNote Password Reset",
		},
	}
	meta := filepath.Join("templates", "meta", "auth.html.tmpl")
	body := filepath.Join("templates", "reset.html.tmpl")
	if r.Method == "POST" {
		r.ParseForm()
		data.Heading, data.Message = "Success!", "A password reset link has been sent to your email address."
		meta = filepath.Join("templates", "meta", "info.html.tmpl")
		body = filepath.Join("templates", "info.html.tmpl")

		if err := models.TryReset(r.Form["mail"][0]); err != nil {
			Handle(err)
			displayHTTPError(w, r, http.StatusInternalServerError)
			return
		}
	}

	tmpl, err := template.ParseFiles(base, meta, body)
	Handle(err)
	tmpl.ExecuteTemplate(w, "layout", data)
	return
}
