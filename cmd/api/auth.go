package main

import (
	"net/http"

	"github.com/harshvse/go-api/internal/store"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required, max=100"`
	Email    string `json:"email" validate:"required, email, max=255"`
	Password string `json:"password" validate:"required, min=3, max=72"`
}

// RegisterUser godoc
//
//	@Summary		Register a new user
//	@Description	Register a user with username, email and password
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User Credentials"
//	@Success		201		{object}	store.User			"User Registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/authentication/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload

	if err := readJson(w, r, payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	// Hash the password
	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// store the user

	if err := app.jsonResponse(w, http.StatusCreated, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}
