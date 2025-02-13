package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/harshvse/go-api/docs"
	"github.com/harshvse/go-api/internal/auth"
	"github.com/harshvse/go-api/internal/mailer"
	"github.com/harshvse/go-api/internal/store"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type application struct {
	config        config
	store         store.Storage
	logger        *zap.SugaredLogger
	mailer        mailer.Client
	authenticator auth.Authenticator
}

type config struct {
	addr        string
	auth        authConfig
	db          dbConfig
	mail        mailConfig
	env         string
	version     string
	frontendURL string
}

type authConfig struct {
	basic basicConfig
}
type basicConfig struct {
	username string
	password string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type mailConfig struct {
	exp        time.Duration
	fromEmail  string
	mailTrap   mailTrap
	maxRetries int
}

type mailTrap struct {
	apikey string
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
		r.With(app.BasicAuthMiddleware()).
			Get("/health", app.healthCheckHandler)

		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", app.registerUserHandler)
		})

		// posts
		r.Route("/posts", func(r chi.Router) {
			r.Post("/create", app.createNewPostHandler)
			r.Route("/{postId}", func(r chi.Router) {
				r.Use(app.postContextMiddleware)
				r.Get("/", app.getPostHandler)
				r.Delete("/", app.deletePostHandler)
				r.Patch("/", app.updatePostHandler)
			})
		})

		// comments
		r.Route("/comments", func(r chi.Router) {
			r.Post("/create", app.createCommentHandler)
			r.Get("/{postId}", app.getCommentByPostIDHandler)
		})

		// users
		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", app.activateUserHandler)
			r.Route("/{userId}", func(r chi.Router) {
				r.Use(app.userContextMiddleware)
				r.Get("/", app.getUserByIDHandler)
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unFollowUserHandler)
			})
			r.Group(func(r chi.Router) {
				r.Get("/feed", app.getUserFeedHandler)
			})
		})
	})

	return r
}

// the run function takes a mux which is responsible for routing and deploys the code
func (app *application) run(mux http.Handler) error {
	docs.SwaggerInfo.Version = app.config.version
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/v1"
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute * 2,
	}

	app.logger.Infow("Server started", "addr", app.config.addr, "env", app.config.env)

	return srv.ListenAndServe()
}
