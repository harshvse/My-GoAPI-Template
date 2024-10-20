package main

import (
	"log"

	"github.com/harshvse/go-api/internal/db"
	"github.com/harshvse/go-api/internal/env"
	"github.com/harshvse/go-api/internal/store"
	"github.com/joho/godotenv"
)

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
		env:     env.GetString("ENVIRONMENT", "Development"),
		version: env.GetString("APIVERSION", "v1"),
	}

	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Panic("Database connection failed ", err)
	}
	defer db.Close()
	log.Println("Database connected")

	store := store.NewPostgresStorage(db)

	// inject dependencies into the server
	app := &application{
		config: cfg,
		store:  store,
	}

	// load all the routes
	mux := app.mount()

	// run the server
	log.Fatal(app.run(mux))
}
