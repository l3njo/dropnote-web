package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type signupData struct {
	Name    string `json:"name"`
	Mail    string `json:"mail"`
	Pass    string `json:"pass"`
	confirm string
}

type loginData struct {
	Mail string `json:"mail"`
	Pass string `json:"pass"`
}

func (s *signupData) validate() (bool, string) {
	if s.Pass != s.confirm {
		return false, "The passwords did not match."
	}
	return true, ""
}

func (s *signupData) tryAuth() (*SessionData, bool) {
	var session *SessionData
	url := fmt.Sprintf("%suser/new", api)
	requestBody, err := json.Marshal(s)
	if err != nil {
		return session, false
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return session, false
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	if status := result["status"].(bool); !status {
		return session, false
	}

	session.buildSession(result)
	return session, true
}

func (l *loginData) tryAuth() (*SessionData, bool) {
	var session *SessionData
	url := fmt.Sprintf("%suser/login", api)
	requestBody, err := json.Marshal(l)
	if err != nil {
		return session, false
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return session, false
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	if status := result["status"].(bool); !status {
		return session, false
	}

	session.buildSession(result)
	return session, true
}

func tryReset(mail string) bool {
	url := fmt.Sprintf("%suser/action/reset", api)
	requestBody, err := json.Marshal(map[string]string{"mail": mail})
	if err != nil {
		return false
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	if status := result["status"].(bool); !status {
		return false
	}

	return true
}
