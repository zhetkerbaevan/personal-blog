package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
)

var Validator = validator.New()

func ParseJSON(r *http.Request, payload any) error {
	//check if body is empty
	if r.Body == nil {
		return fmt.Errorf("MISSING REQUEST BODY")
	}
	return json.NewDecoder(r.Body).Decode(payload) //create new decoder which will read data from r.Body, then decode reads data from json and decode it to payload
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	//add header (json) for response
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v) //create new encoder which will write data to w, then encode will write data from v to json
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}