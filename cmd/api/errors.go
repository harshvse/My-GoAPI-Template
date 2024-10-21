package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("[-] internal server error: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())
	writeJsonError(w, http.StatusInternalServerError, "server encountered an error")
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("[-] bad request error: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())
	writeJsonError(w, http.StatusBadRequest, "bad request")
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("[-] not found error: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())
	writeJsonError(w, http.StatusNotFound, "resource not found!")
}
