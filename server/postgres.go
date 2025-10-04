package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

var pgx_db *pgx.Conn

func initPostgres() {

	connString := "postgres://admin:admin123@localhost:5432/orders"

	var err error
	pgx_db, err = pgx.Connect(context.Background(), connString)

	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	fmt.Println("Connected to PostgreSQL database!")
	_ = pgx_db
}
