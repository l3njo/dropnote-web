package controllers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// SignupHandler handles the "/signup" route.
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionCookie)
	Handle(err)
	isAuth := checkAuth(session)
	data := info{Title: "Sign Up"}
	meta := filepath.Join("templates", "meta", "auth.html.tmpl")
	menu := getMenu(isAuth)
	body := filepath.Join("templates", "signup.html.tmpl")
	if r.Method == "POST" {
		r.ParseForm()
		data.Heading = "Error!"
		meta = filepath.Join("templates", "meta", "info.html.tmpl")
		body = filepath.Join("templates", "info.html.tmpl")
		payload := signupData{
			Name:    r.Form["name"][0],
			Mail:    r.Form["mail"][0],
			Pass:    r.Form["pass"][0],
			confirm: r.Form["confirm"][0],
		}
		if ok, msg := payload.validate(); !ok {
			data.Message = msg
		}

		if sData, ok := payload.tryAuth(); ok {
			session.Values["isAuth"] = true
			session.Values["data"] = sData
			log.Println(session.Save(r, w))
			http.Redirect(w, r, getNext(r), http.StatusFound)
			return
		}
		data.Message = "Your signup credentials are invalid."
	}

	tmpl, err := template.ParseFiles(base, menu, meta, body)
	Handle(err)
	tmpl.ExecuteTemplate(w, "layout", data)
	return
}

// LoginHandler handles the "/login" route.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionCookie)
	Handle(err)
	isAuth := checkAuth(session)
	data := info{Title: "Log In"}
	meta := filepath.Join("templates", "meta", "auth.html.tmpl")
	menu := getMenu(isAuth)
	body := filepath.Join("templates", "login.html.tmpl")
	if r.Method == "POST" {
		r.ParseForm()
		data.Heading = "Error!"
		meta = filepath.Join("templates", "meta", "info.html.tmpl")
		body = filepath.Join("templates", "info.html.tmpl")
		payload := loginData{
			Mail: r.Form["mail"][0],
			Pass: r.Form["pass"][0],
		}

		if sData, ok := payload.tryAuth(); ok {
			session.Values["isAuth"] = true
			session.Values["data"] = sData
			log.Println(session.Save(r, w))
			http.Redirect(w, r, getNext(r), http.StatusFound)
			return
		}
		data.Message = "Your login credentials are invalid."
	}

	tmpl, err := template.ParseFiles(base, menu, meta, body)
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
	log.Println(session.Save(r, w))
	http.Redirect(w, r, getNext(r), http.StatusFound)
	return
}

// ResetHandler handles the "/reset" route.
func ResetHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionCookie)
	Handle(err)
	isAuth := checkAuth(session)
	data := info{Title: "Reset Password"}
	meta := filepath.Join("templates", "meta", "auth.html.tmpl")
	menu := getMenu(isAuth)
	body := filepath.Join("templates", "reset.html.tmpl")
	if r.Method == "POST" {
		r.ParseForm()
		data.Heading, data.Message = "Success!", "A password reset link has been sent to your email address."
		meta = filepath.Join("templates", "meta", "info.html.tmpl")
		body = filepath.Join("templates", "info.html.tmpl")

		if err := tryReset(r.Form["mail"][0]); err != nil {
			Handle(err)
			data.Heading, data.Message = "Error!", "Something has gone horribly wrong."
		}
	}

	tmpl, err := template.ParseFiles(base, menu, meta, body)
	Handle(err)
	tmpl.ExecuteTemplate(w, "layout", data)
	return
}
