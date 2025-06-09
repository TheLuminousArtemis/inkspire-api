package main

import (
	"log"
	"time"

	"github.com/theluminousartemis/socialnews/internal/db"
	"github.com/theluminousartemis/socialnews/internal/env"
	"github.com/theluminousartemis/socialnews/internal/mailer"
	"github.com/theluminousartemis/socialnews/internal/store"
	"go.uber.org/zap"
)

const version = "0.0.1"

//	@title			Wise.ly
//	@description	API for Wise.ly, a community driven Q&A platform.

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
//
// @securityDefinitions.apikey	ApiAuthKey
// @in							header
// @name						Authorization
// @description
func main() {
	cfg := &Config{
		addr:   env.GetString("ADDR", ":8080"),
		apiURL: env.GetString("API_URL", "localhost:3000"),

		db: &dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5432/wisely?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "production"),
		mail: mailConfig{
			exp:       time.Hour * 24,
			fromEmail: env.GetString("FROM_EMAIL", ""),
			username:  env.GetString("MAILTRAP_USERNAME", ""),
			password:  env.GetString("MAILTRAP_PASSWORD", ""),
		},
		auth: authConfig{
			basic: basicConfig{
				user: env.GetString("AUTH_BASIC_USER", "admin"),
				pass: env.GetString("AUTH_BASIC_PASS", "admin"),
			},
		},
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:4000"),
	}
	//logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	//database
	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	logger.Info("database connection pool established")
	store := store.NewPostgresStorage(db)

	//mailer
	// mailer := mailer.NewSendgrid(cfg.mail.sendGrid.apiKey, cfg.mail.fromEmail)
	mailtrap, err := mailer.NewMailTrapClient(cfg.mail.username, cfg.mail.fromEmail, cfg.mail.password)
	if err != nil {
		log.Fatal(err)
	}

	app := &application{
		config:  cfg,
		storage: store,
		l:       logger,
		mailer:  mailtrap,
	}

	mux := app.Mount()
	app.Start(mux)
}
