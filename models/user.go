package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// User represents a user
type User struct {
	User    string `json:"ID,omitempty"`
	Name    string `json:"name,omitempty"`
	Mail    string `json:"mail,omitempty"`
	Pass    string `json:"pass,omitempty"`
	Auth    string `json:"auth,omitempty"`
	Confirm string `json:"-"`
}

type userResult struct {
	User    `json:"data"`
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

// ValidateSignup checks signup details
func (u *User) ValidateSignup() error {
	if u.Pass != u.Confirm {
		return errors.New("the passwords did not match")
	}
	return nil
}

// TrySignup attempts a signup
func (u *User) TrySignup() error {
	requestBody, err := json.Marshal(u)
	if err != nil {
		return err
	}
	resp, err := http.Post(fmt.Sprintf("%susers/new", api), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	result := &userResult{}
	json.NewDecoder(resp.Body).Decode(result)
	if !result.Status {
		return errors.New(result.Message)
	}

	*u = result.User
	return nil
}

// TryLogin attempts a login
func (u *User) TryLogin() error {
	requestBody, err := json.Marshal(u)
	if err != nil {
		return err
	}

	resp, err := http.Post(fmt.Sprintf("%susers/login", api), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	result := &userResult{}
	json.NewDecoder(resp.Body).Decode(result)
	if !result.Status {
		return errors.New(result.Message)
	}

	*u = result.User
	return nil
}

// TryDelete attempts to delete a user
func (u *User) TryDelete() error {
	request, err := http.NewRequest("DELETE", fmt.Sprintf("%sme/delete", api), nil)
	if err != nil {
		return err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", u.Auth))
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	result := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&result)
	if status := result["status"].(bool); !status {
		return errors.New(result["message"].(string))
	}

	return nil
}

// TryReset attempts a password reset
func TryReset(mail string) error {
	requestBody, err := json.Marshal(map[string]string{"mail": mail})
	if err != nil {
		return err
	}

	resp, err := http.Post(fmt.Sprintf("%susers/actions/reset", api), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	result := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&result)
	if status := result["status"].(bool); !status {
		return errors.New(result["message"].(string))
	}

	return nil
}
