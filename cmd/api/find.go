package main

import (
	"encoding/json"
	"net/http"
)

type Request struct {
	Profession 	string `json:"profession,omitempty"`
	JobTitle 	string `json:"jobTitle,omitempty"`
	CompanyName string `json:"companyName,omitempty"`
	EmailAddr 	string `json:"email,omitempty"`
}

func (app *application) handleFindJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}

	var input Request

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		app.logger.Fatal(err)
	}

	// Send an acknowledgement as a response
	// Send an internal message to the Manager
	// Start the scraping process 
	// Send an email to the addr after scraping
}