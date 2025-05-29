package main

import (
	"log"

	"github.com/theluminousartemis/socialnews/internal/db"
	"github.com/theluminousartemis/socialnews/internal/env"
	"github.com/theluminousartemis/socialnews/internal/store"
)

const version = "0.0.1"

func main() {
	cfg := &Config{
		addr: env.GetString("ADDR", ":8080"),
		db: &dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5432/socialnews?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
	}
	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	log.Printf("Database connected at %s", cfg.addr)
	store := store.NewPostgresStorage(db)
	app := &application{
		config:  cfg,
		storage: store,
	}
	mux := app.Mount()
	app.Start(mux)
}
