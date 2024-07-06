package main

import (
	"encoding/json"
	"net/http"
)

func (app *Application) errorResponse(w http.ResponseWriter, errMsg string) {
	errorMsg := struct{
		Error string		
	}{
		Error: errMsg,
	}
	// fmt.Errorf for making errors
	app.Logger.Printf("Error: %s\n", errMsg)
	writeJSON(w, errorMsg)
}


func writeJSON(w http.ResponseWriter, data interface{}) error {
	// w.WriteHeader(http.StatusInternalServerError) -> add status codes later
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Send back error resp
		return err
	}

	return nil
}