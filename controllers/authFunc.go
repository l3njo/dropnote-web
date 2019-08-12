package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
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
	session := &SessionData{}
	url := fmt.Sprintf("%susers/new", api)
	requestBody, err := json.Marshal(s)
	if err != nil {
		return session, false
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return session, false
	}
	defer resp.Body.Close()

	result := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&result)
	if status := result["status"].(bool); !status {
		return session, false
	}
	log.Println(result) // DEBUG

	session.buildSession(result)
	return session, true
}

func (l *loginData) tryAuth() (*SessionData, bool) {
	session := &SessionData{}
	url := fmt.Sprintf("%susers/login", api)
	requestBody, err := json.Marshal(l)
	if err != nil {
		return session, false
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return session, false
	}
	defer resp.Body.Close()

	result := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&result)
	if status := result["status"].(bool); !status {
		return session, false
	}

	session.buildSession(result)
	return session, true
}

func tryReset(mail string) error {
	url := fmt.Sprintf("%susers/actions/reset", api)
	requestBody, err := json.Marshal(map[string]string{"mail": mail})
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	result := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&result)
	if status := result["status"].(bool); !status {
		return errors.New(result["message"])
	}

	return nil
}
