package main

import (
	"log"

	"github.com/harshvse/go-api/internal/db"
	"github.com/harshvse/go-api/internal/store"
)

func main() {

	addr := "postgres://user:password@localhost:5432/golangdb?sslmode=disable"
	maxOpenConns := 20
	maxIdleConns := 20
	maxIdleTime := "15m"

	conn, err := db.New(addr, maxOpenConns, maxIdleConns, maxIdleTime)
	if err != nil {
		log.Panic("failed to connected to database", err)
	}
	defer conn.Close()
	log.Println("database connected")

	store := store.NewPostgresStorage(conn)
	db.Seed(store, conn)
}
