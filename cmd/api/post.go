package main

import (
	"net/http"

	"github.com/harshvse/go-api/internal/store"
)

type CreatePostPayload struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (app *application) createNewPost(w http.ResponseWriter, r *http.Request) {
	var postPayload CreatePostPayload

	if err := readJson(w, r, &postPayload); err != nil {
		writeJsonError(w, http.StatusBadRequest, err.Error())
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
		writeJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeJson(w, http.StatusCreated, post); err != nil {
		writeJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
