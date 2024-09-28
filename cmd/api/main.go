package main

import (
	"log"

	"github.com/harshvse/go-api/internal/env"
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
	}

	// inject dependencies into the server
	app := &application{
		config: cfg,
	}

	// load all the routes
	mux := app.mount()

	// run the server
	log.Fatal(app.run(mux))
}
