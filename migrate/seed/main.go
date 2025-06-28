package main

import (
	"log"

	"github.com/theluminousartemis/inkspire/internal/db"
	"github.com/theluminousartemis/inkspire/internal/env"
	"github.com/theluminousartemis/inkspire/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5432/wisely?sslmode=disable")
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	store := store.NewPostgresStorage(conn)
	db.Seed(store, conn)
}
