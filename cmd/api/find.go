package main

import (
	"encoding/json"
	"net/http"

	actors "github.com/kerosiinikone/go-actors-project/internal"
)

type Request struct {
	Profession 	string `json:"profession,omitempty"`
	JobTitle 	string `json:"jobTitle,omitempty"`
	CompanyName string `json:"companyName,omitempty"`
	EmailAddr 	string `json:"email,omitempty"`
}

func (app *Application) handleFindJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}

	var input Request

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		app.Logger.Fatal(err)
	}

	// Validate the response !!!

	// Send an internal message to the Manager
	actors.HandleNewJobRequest()
	
	// Start the scraping process 
	// Send an email to the addr after scraping
	
	// Send an acknowledgement as a response
	writeJSON(w, map[string]bool{
		"success": true,
	}, 201)

}