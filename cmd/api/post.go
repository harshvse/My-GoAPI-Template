package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/harshvse/go-api/internal/store"
)

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=10000"`
	Tags    []string `json:"tags"`
}

func (app *application) createNewPostHandler(w http.ResponseWriter, r *http.Request) {
	var postPayload CreatePostPayload

	if err := readJson(w, r, &postPayload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(postPayload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()

	// TODO: Change to real userid
	userId := 1

	post := &store.Post{
		Title:   postPayload.Title,
		Content: postPayload.Content,
		Tags:    postPayload.Tags,
		UserID:  int64(userId),
	}

	if err := app.store.Posts.Create(ctx, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJson(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement get all posts for a user ID
	postIdString := chi.URLParam(r, "postId")
	postId, err := strconv.ParseInt(postIdString, 10, 64)

	if err != nil {
		app.badRequestError(w, r, err)
		return
	}
	ctx := r.Context()

	post, err := app.store.Posts.GetByID(ctx, postId)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := writeJson(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
