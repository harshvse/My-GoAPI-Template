package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "OK!",
		"env":     app.config.env,
		"version": app.config.version,
	}
	if err := writeJson(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
