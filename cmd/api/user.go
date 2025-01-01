package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/harshvse/go-api/internal/store"
)

type userKey string

const userctx userKey = "user"

// GetUser godoc
//
//	@Summary		Fetch a user profile
//	@Description	Fetch a user profile by ID
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	store.User
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{id} [get]
func (app *application) getUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)
	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type FollowUser struct {
	UserID int64 `json:"user_id"`
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followedUser := getUserFromCtx(r)

	var payload FollowUser
	if err := readJson(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	ctx := r.Context()
	if err := app.store.Followers.Follow(ctx, followedUser.ID, payload.UserID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) unFollowUserHandler(w http.ResponseWriter, r *http.Request) {
	unFollowedUser := getUserFromCtx(r)

	var payload FollowUser
	if err := readJson(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	ctx := r.Context()
	if err := app.store.Followers.UnFollow(ctx, unFollowedUser.ID, payload.UserID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// ActivateUser godoc
//
//	@Summary		Activate a new user from email
//	@Description	Activate a user from the email invitation sent during register
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			payload	path		string	true	"Invitation Token"
//	@Success		204		{String}	string	"User account activated"
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/activate/{token} [put]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	err := app.store.Users.ActivateUser(r.Context(), token)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.badRequestError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}
	if err := app.jsonResponse(w, http.StatusNoContent, ""); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIdString := chi.URLParam(r, "userId")
		userId, err := strconv.ParseInt(userIdString, 10, 64)
		if err != nil {
			app.badRequestError(w, r, err)
			return
		}

		ctx := r.Context()
		user, err := app.store.Users.GetByID(ctx, userId)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundError(w, r, err)
				return
			default:
				app.internalServerError(w, r, err)
				return
			}
		}
		ctx = context.WithValue(ctx, userctx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromCtx(r *http.Request) *store.User {
	user, _ := r.Context().Value(userctx).(*store.User)
	return user
}
