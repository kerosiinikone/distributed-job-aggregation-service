package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestFind(t *testing.T) {
	b := new(bytes.Buffer)
	body := map[string]string{
		"location": "Helsinki",
		"email": "rand@gmail.com",
	}

	json.NewEncoder(b).Encode(body)
	resp, err := http.Post("http://localhost:3000/find", "application/json", b)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	var output struct {
		Success bool `json:"success"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		t.Fatal(err)
	}

	if ok := output.Success; !ok {
		t.Fatal(err)
	}
}