package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Note represents entire note as provided by API
type Note struct {
	Voucher string `json:"ID"`
	Subject string `json:"subject"`
	Content string `json:"content"`
	Dropped string `json:"created_at"`
	Visible bool   `json:"visible"`
}

// NoteSlice represents an array of notes
type NoteSlice struct {
	Notes  []Note `json:"data"`
	Status bool   `json:"status"`
}

func getNotes(auth string) ([]Note, error) {
	notes := []Note{}
	url := fmt.Sprintf("%sme/notes", api)
	auth = fmt.Sprintf("Bearer %s", auth)
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
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		log.Println(err)
	}
	if result.Status {
		resultData := result.Notes
		for _, v := range resultData {
			if len(v.Subject) > 10 {
				v.Subject = v.Subject[:9] + "…"
			}
			if len(v.Content) > 20 {
				v.Content = v.Content[:19] + "…"
			}
			now, err := time.Parse("2006-02-01T15:04:05Z", v.Dropped)
			if err != nil {
				return nil, err
			}
			v.Dropped = now.Format("2006-02-01 15:04:05")
			notes = append(notes, v)
		}
	}

	return notes, nil
}
