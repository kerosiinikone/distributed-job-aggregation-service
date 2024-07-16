package main

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
)

var (
	EmailRegEx = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$") 
)

type JobRequest struct {
	Keywords  	[]string `json:"keywords,omitempty"`
	EmailAddr   string `json:"email,omitempty"`
	Location	string `json:"location,omitempty"`
	maxPages	int
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

	if input.EmailAddr == "" || !EmailRegEx.Match([]byte(input.EmailAddr)) {
		app.errorResponse(w, "Invalid Email Address")
		return
	}

	// Why not allow actors to access the main app struct?
	input.maxPages = app.Cfg.MaxPages

	// Send an internal message to the Manager
	err := app.handleNewJobRequest(&input)
	if err != nil {
		app.errorResponse(w, err.Error())
		return
	}
	
	// Send an acknowledgement as a response
	app.writeJSON(w, map[string]bool{
		"success": true,
	})

}

func (app *Application) handleNewJobRequest(input *JobRequest) error {
	// Send Message to Manager
	// Some additional logic / validation if necessary
	app.Engine.Send(app.MPid, input)
	return nil
}
