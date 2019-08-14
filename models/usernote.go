package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// NoteSlice represents an array of notes
type NoteSlice struct {
	Notes   []Note `json:"data"`
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

// GetNotes returns an array of notes created by provided user
func (u *User) GetNotes() ([]Note, error) {
	notes := []Note{}
	url := fmt.Sprintf("%sme/notes", api)
	auth := fmt.Sprintf("Bearer %s", u.Auth)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", auth)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result := &NoteSlice{}
	if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, err
	}

	if !result.Status {
		return nil, errors.New(result.Message)
	}

	for _, v := range result.Notes {
		if len(v.Subject) > 10 {
			v.Subject = v.Subject[:9] + "…"
		}
		if len(v.Content) > 20 {
			v.Content = v.Content[:19] + "…"
		}
		now, err := time.Parse("2006-01-02T15:04:05Z", v.Dropped)
		if err != nil {
			return nil, err
		}
		v.Dropped = now.Format("2006-01-02 15:04:05")
		notes = append(notes, v)
	}

	return notes, nil
}
