package main

import (
	"log"

	"github.com/theluminousartemis/socialnews/internal/db"
	"github.com/theluminousartemis/socialnews/internal/env"
	"github.com/theluminousartemis/socialnews/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5432/socialnews?sslmode=disable")
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	store := store.NewPostgresStorage(conn)
	db.Seed(store, conn)
}
