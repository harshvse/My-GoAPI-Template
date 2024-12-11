package main

import (
	"net/http"

	"github.com/harshvse/go-api/internal/store"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	//pagination,filters
<<<<<<< HEAD
	fq := store.PaginatedFeedQuery{
		Limit:  20,
=======

	fq := store.PaginatedFeedQuery{
		Limit:  10,
>>>>>>> 41df07bd1f187bca45379cfaea95304f598a167e
		Offset: 0,
		Sort:   "desc",
	}

	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestError(w, r, err)
<<<<<<< HEAD
=======
		return
>>>>>>> 41df07bd1f187bca45379cfaea95304f598a167e
	}

	if err := Validate.Struct(fq); err != nil {
		app.badRequestError(w, r, err)
<<<<<<< HEAD
	}
	ctx := r.Context()
	feed, err := app.store.Posts.GetUserFeed(ctx, int64(3), fq)
=======
		return
	}

	ctx := r.Context()
	feed, err := app.store.Posts.GetUserFeed(ctx, int64(2), fq)
>>>>>>> 41df07bd1f187bca45379cfaea95304f598a167e

	if err != nil {
		app.internalServerError(w, r, err)
	}
	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
	}
}
