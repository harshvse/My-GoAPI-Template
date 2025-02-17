package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/harshvse/go-api/internal/mailer"
	"github.com/harshvse/go-api/internal/store"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}
type UserWithToken struct {
	*store.User
	Token string `json:"token"`
}

// RegisterUser godoc
//
//	@Summary		Register a new user
//	@Description	Register a user with username, email and password
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User Credentials"
//	@Success		201		{object}	UserWithToken		"User Registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/authentication/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload

	if err := readJson(w, r, &payload); err != nil {
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

	ctx := r.Context()

	// Generate User token
	plainToken := uuid.New().String()

	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	// store the user
	err := app.store.Users.CreateAndInviteUser(ctx, user, hashToken, app.config.mail.exp)
	if err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			app.badRequestError(w, r, err)
		case store.ErrDuplicateUsername:
			app.badRequestError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}
	userWithToken := UserWithToken{
		User:  user,
		Token: plainToken,
	}

	// TODO update the addr to frontend url
	activationUrl := fmt.Sprintf("%s/confirm/%s", app.config.addr, plainToken)
	isProdEnv := app.config.env == "production"
	vars := struct {
		Username      string
		ActivationUrl string
	}{
		Username:      user.Username,
		ActivationUrl: activationUrl,
	}

	// Send activate email
	for i := 0; i < app.config.mail.maxRetries; i++ {
		_, err = app.mailer.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, vars, !isProdEnv)
		if err != nil {
			app.logger.Errorw("error sending the activation url", err)
			if i >= app.config.mail.maxRetries-1 {
				app.logger.Info("max retries for email send reached for user %v. Aborting user creation.", user.Username)
				// roll back user creation if email fails
				if err := app.store.Users.Delete(ctx, user.ID); err != nil {
					app.logger.Errorw("failed to delte the user", err)
					return
				}
			}
			app.internalServerError(w, r, err)
			continue
		}
		break
	}
	if err := app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		app.internalServerError(w, r, err)
	}
}

type CreateUserTokenPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

// CreateTokenHandler godoc
//
//	@Summary		Creates a new token
//	@Description	Create a new token for a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateUserTokenPayload	true	"User Credentials"
//	@Success		201		{string}	string					"Token"
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/token [post]
func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Parse payload credentials
	var payload CreateUserTokenPayload

	if err := readJson(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// fetch the user if he exists from the credentials
	user, err := app.store.Users.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			//TODO make this not found but don't return to frontend
			app.unauthorizedError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}
	// generate the token -> add claims

	// send it to the client
	if err := app.jsonResponse(w, http.StatusCreated, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}
