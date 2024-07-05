package main

import (
	"encoding/json"
	"net/http"
)


func (app *Application) findJobHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.errorResponse(w, "Wrong Method", 500)
		return
	}

	var input JobRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		app.errorResponse(w, err.Error(), 500)

	}

	// Validate the response !!!

	// Send an internal message to the Manager
	err := app.handleNewJobRequest(&input)
	if err != nil {
		app.errorResponse(w, err.Error(), 500)
		return
	}
	
	// Send an acknowledgement as a response
	writeJSON(w, map[string]bool{
		"success": true,
	}, 201)

}

func (app *Application) handleNewJobRequest(input *JobRequest) error {
	// Send Message to Manager
	// Some additional logic / validation if necessary
	app.Engine.Send(app.MPid, input)

	return nil
}
