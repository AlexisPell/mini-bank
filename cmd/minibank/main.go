package main

import (
	"database/sql"
	db "github.com/alexispell/minibank/db/sqlc"
	"github.com/alexispell/minibank/internal/api"
	_ "github.com/lib/pq"
	"log"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://user:123456@localhost:5432/minibank?sslmode=disable"
	serverAddress = "127.0.0.1:8080"
)

func main() {
	// connect to db
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to database", err)
	}

	// create db store
	store := db.NewStore(conn)

	// launch the server
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("Cannot start server", err)
	}
}
