package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, req *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": app.config.version,
	}
	if err := writeJson(w, http.StatusOK, data); err != nil {
		writeJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// CRUD CREATE READ UPDATE DELETE
