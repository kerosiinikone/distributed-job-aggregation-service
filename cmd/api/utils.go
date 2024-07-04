package main

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, data interface{}, statusCode int) error {
	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Send back error resp
		return err
	}

	return nil
}