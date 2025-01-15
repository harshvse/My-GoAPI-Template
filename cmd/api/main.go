package main

import (
	"log"
	"time"

	"github.com/harshvse/go-api/internal/db"
	"github.com/harshvse/go-api/internal/env"
	"github.com/harshvse/go-api/internal/mailer"
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
		mail: mailConfig{
			exp:        time.Hour * 24 * 3, // 3 days
			fromEmail:  env.GetString("SENDGRID_EMAIL", "hello@demomailtrap.com"),
			maxRetries: 5,
			mailTrap: mailTrap{
				apikey: env.GetString("MAILTRAP_API_KEY", "bad-key"),
			},
		},
		env:         env.GetString("ENVIRONMENT", "DEVELOPMENT"),
		version:     env.GetString("APIVERSION", "UNDEFINED"),
		frontendURL: env.GetString("Frontend_URL", "http://localhost:3000"),
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

	// Mailer
	mailer, err := mailer.NewMailTrapClient(cfg.mail.mailTrap.apikey, cfg.mail.fromEmail)
	if err != nil {
		logger.Errorw("mailer creation failed", err)
		return
	}

	// inject dependencies into the server
	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
		mailer: mailer,
	}

	// load all the routes
	mux := app.mount()

	// run the server
	logger.Fatal(app.run(mux))
}
