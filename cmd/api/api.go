package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/harshvse/go-api/internal/store"
)

type application struct {
	config config
	store  store.Storage
}

type config struct {
	addr    string
	db      dbConfig
	env     string
	version string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

// this is where all the middlewares and the routes will be handled
func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		r.Route("/posts", func(r chi.Router) {
			r.Post("/create", app.createNewPostHandler)
			r.Route("/{postId}", func(r chi.Router) {
				r.Get("/", app.getPostHandler)
			})
		})

		r.Route("/comments", func(r chi.Router) {
			r.Post("/create", app.createCommentHandler)
			r.Get("/{postId}", app.getCommentByPostIDHandler)
		})

	})

	return r
}

// the run function takes a mux which is responsible for routing and deploys the code
func (app *application) run(mux http.Handler) error {

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute * 2,
	}

	log.Printf("Server started on %s", app.config.addr)

	return srv.ListenAndServe()
}
