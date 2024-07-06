package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type JobRequest struct {
	Keywords  	[]string `json:"keywords,omitempty"`
	EmailAddr   string `json:"email,omitempty"`
	Location	string `json:"location,omitempty"`
}

func (app *Application) findJobHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.errorResponse(w, "Wrong Method")
		return
	}

	var input JobRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		if err != io.EOF {
			app.errorResponse(w, err.Error())
			return
		}
	}

	// Validate the response !!!
	// Send an internal message to the Manager
	err := app.handleNewJobRequest(&input)
	if err != nil {
		app.errorResponse(w, err.Error())
		return
	}
	
	// Send an acknowledgement as a response
	writeJSON(w, map[string]bool{
		"success": true,
	})

}

func (app *Application) handleNewJobRequest(input *JobRequest) error {
	// Send Message to Manager
	// Some additional logic / validation if necessary
	app.Engine.Send(app.MPid, input)
	return nil
}
