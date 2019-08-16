package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

// Note represents entire note as provided by API
type Note struct {
	Voucher string `json:"ID,omitempty"`
	Subject string `json:"subject,omitempty"`
	Content string `json:"content,omitempty"`
	Dropped string `json:"created_at,omitempty"`
	Creator string `json:"creator,omitempty"`
	Visible bool   `json:"visible,omitempty"`
}

type noteResult struct {
	Note    `json:"data"`
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

// ValidateGet checks validity of provided code
func (n *Note) ValidateGet() error {
	if _, e := uuid.FromString(n.Voucher); e != nil {
		return errors.New("invalid code")
	}
	return nil
}

// Get returns populates the provided note with data
func (n *Note) Get(auth string) error {
	url := fmt.Sprintf("%snotes/%s", api, n.Voucher)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	if auth != "" {
		auth := fmt.Sprintf("Bearer %s", auth)
		request.Header.Add("Authorization", auth)
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	result := &noteResult{}
	json.NewDecoder(resp.Body).Decode(&result)
	if !result.Status {
		return errors.New(result.Message)
	}

	*n = result.Note
	return nil
}

// ValidatePost checks validity of provided note
func (n *Note) ValidatePost() error {
	if n.Subject == "" {
		return errors.New("empty subject")
	}
	if n.Content == "" {
		return errors.New("empty content")
	}
	return nil
}

// Post saves the provided note
func (n *Note) Post(auth string) error {
	requestBody, err := json.Marshal(n)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%snotes/new", api)
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	request.Header.Set("Content-type", "application/json")
	if auth != "" {
		auth := fmt.Sprintf("Bearer %s", auth)
		request.Header.Add("Authorization", auth)
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	result := &noteResult{}
	json.NewDecoder(resp.Body).Decode(&result)
	if !result.Status {
		return errors.New(result.Message)
	}

	*n = result.Note
	return nil
}

// Toggle changes the provided note's visibility
func (n *Note) Toggle(auth string) error {
	url := fmt.Sprintf("%sme/notes/toggle/%s", api, n.Voucher)
	request, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", auth))
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	result := &noteResult{}
	json.NewDecoder(resp.Body).Decode(&result)
	if !result.Status {
		return errors.New(result.Message)
	}

	*n = result.Note
	return nil
}

// Update updates the provided note
func (n *Note) Update(auth string) error {
	requestBody, err := json.Marshal(n)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%sme/notes/update/%s", api, n.Voucher)
	request, err := http.NewRequest("PUT", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", auth))
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	result := &noteResult{}
	json.NewDecoder(resp.Body).Decode(&result)
	if !result.Status {
		return errors.New(result.Message)
	}

	*n = result.Note
	return nil
}

// Delete deletes the provided note
func (n *Note) Delete(auth string) error {
	url := fmt.Sprintf("%sme/notes/delete/%s", api, n.Voucher)
	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", auth))
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

	*n = Note{}
	return nil
}

// ParseDate converts the date into a more readable format
func (n *Note) ParseDate() error {
	now, err := time.Parse("2006-01-02T15:04:05Z", n.Dropped)
	if err != nil {
		return err
	}
	n.Dropped = now.Format("2006-01-02 15:04:05")
	return nil
}
