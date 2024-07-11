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
	app.writeJSON(w, errorMsg)
}


func (app *Application) writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		errMsg := []byte(err.Error())
		w.Write(errMsg) // For now
	}
}