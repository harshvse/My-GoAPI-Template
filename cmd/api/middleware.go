package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// read the auth header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthorizedBasicError(w, r, fmt.Errorf("no Authorization header found"))
				return
			}
			// parse and convert it from base64
			parts := strings.Split(authHeader, " ")
			if len(parts) < 2 || parts[0] != "Basic" {
				app.unauthorizedBasicError(w, r, fmt.Errorf("invalid Authorization header found"))
				return
			}

			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unauthorizedBasicError(w, r, err)
				return
			}

			username := app.config.auth.basic.username
			password := app.config.auth.basic.password

			credentials := strings.SplitN(string(decoded), ":", 2)

			if len(credentials) != 2 || credentials[0] != username || credentials[1] != password {
				app.unauthorizedBasicError(w, r, fmt.Errorf("invalid credentials"))
				return
			}

			// check the credentials against what we have stored

			next.ServeHTTP(w, r)
		})
	}
}
