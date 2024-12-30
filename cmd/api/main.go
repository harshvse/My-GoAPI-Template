package main

import (
	"log"

	"github.com/harshvse/go-api/internal/db"
	"github.com/harshvse/go-api/internal/env"
	"github.com/harshvse/go-api/internal/store"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

//	@title			GoLang WebAPI Template With Chi-Postgres
//	@version		1.0
//	@description	This is a sample server for a social media.
//	@termsOfService	/

//	@contact.name	Harsh Verma
//	@contact.url	http://harshvse.in
//	@contact.email	harshvse@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host						localhost:8080
// @BasePath					/v1
//
// @securityDefinitions.apiKey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	// load the .env file into environment variables so env.go can read them
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// set the configuration of the server
	cfg := config{
		addr: env.GetString("SERVER_URL", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_URL", "postgres://user:password@localhost:5432/golangdb?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env:     env.GetString("ENVIRONMENT", "DEVELOPMENT"),
		version: env.GetString("APIVERSION", "UNDEFINED"),
	}

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Database
	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Fatal("Database connection failed ", err)
	}
	defer db.Close()
	logger.Info("Database connected")

	store := store.NewPostgresStorage(db)

	// inject dependencies into the server
	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
	}

	// load all the routes
	mux := app.mount()

	// run the server
	logger.Fatal(app.run(mux))
}
