package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("internal error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJsonError(w, http.StatusInternalServerError, "server encountered an error")
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("bad request", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJsonError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("not found", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJsonError(w, http.StatusNotFound, "resource not found!")
}
