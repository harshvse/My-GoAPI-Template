package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

func writeJson(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json") //set the content type header so the browser will know it's json
	w.WriteHeader(status)                              // set the status header
	return json.NewEncoder(w).Encode(data)             // encode whatever data to json
}

func readJson(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_578                                    // 1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes)) // disallow any request with more than 1MB

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(data)
}

func writeJsonError(w http.ResponseWriter, status int, message string) error {
	type envelop struct {
		Error string `json:"error"`
	}
	return writeJson(w, status, &envelop{Error: message})
}

func (app *application) jsonResponse(w http.ResponseWriter, status int, data any) error {
	type envelop struct {
		Data any `json:"data"`
	}
	return writeJson(w, status, &envelop{data})
}
