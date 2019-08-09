package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

type note struct{
	subject, content string
}

func validate(voucher string) (isValid bool) {
	isValid = true
	if _, e := uuid.FromString(voucher); e != nil{
		isValid = false
	}
	return
}

func getNote(url string) (note, error) {
	data := note{}
	resp, err := http.Get(url)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	resultData := result["data"].(map[string]interface{})
	data.subject = resultData["subject"].(string)
	data.content = resultData["content"].(string)
	return data, nil
}

func postNote(url string, n note) (bool, string) {
	var voucher string
	// TODO Remove map
	requestBody, err := json.Marshal(map[string]string{
		"subject": n.subject,
		"content": n.content,
	}) 
	if err != nil {
		return false, voucher
	}

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