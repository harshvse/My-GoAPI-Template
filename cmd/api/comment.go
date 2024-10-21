package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/harshvse/go-api/internal/store"
)

type CommentPayload struct {
	// UNSAFE make sure to only get userId through JWT in the furutre
	PostID  int64  `json:"post_id"`
	UserID  int64  `json:"user_id"`
	Content string `json:"content"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var commentPayload CommentPayload
	if err := readJson(w, r, &commentPayload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()
	comment := &store.Comment{
		PostID:  int64(commentPayload.PostID),
		UserID:  int64(commentPayload.UserID),
		Content: commentPayload.Content,
	}
	if err := app.store.Comments.Create(ctx, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := writeJson(w, http.StatusOK, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getCommentByPostIDHandler(w http.ResponseWriter, r *http.Request) {
	postIdString := chi.URLParam(r, "postId")
	postId, err := strconv.ParseInt(postIdString, 10, 64)

	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()

	postWithComments, err := app.store.Comments.GetPostByID(ctx, postId)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := writeJson(w, http.StatusOK, postWithComments); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
