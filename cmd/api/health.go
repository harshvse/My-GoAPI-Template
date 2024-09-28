package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Ok!"))
}
