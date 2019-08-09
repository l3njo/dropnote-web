package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

type note struct {
	Subject string `json:"subject"`
	Content string `json:"content"`
}

func validateCode(voucher string) bool {
	if _, e := uuid.FromString(voucher); e != nil {
		return false
	}
	return true
}

func (n *note) getNote(voucher string) error {
	url := fmt.Sprintf("%snote/%s", api, voucher)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	resultData := result["data"].(map[string]interface{})
	n.Subject = resultData["subject"].(string)
	n.Content = resultData["content"].(string)
	return nil
}

func (n *note) postNote() (bool, string) {
	var voucher string
	requestBody, err := json.Marshal(n)
	if err != nil {
		return false, voucher
	}

	url := fmt.Sprintf("%snote/new", api)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return false, voucher
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	status := result["status"].(bool)
	if status {
		resultData := result["data"].(map[string]interface{})
		voucher = resultData["ID"].(string)
	}

	return status, voucher
}
