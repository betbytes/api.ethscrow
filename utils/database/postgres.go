package database

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

var DB *pgxpool.Pool

func ConnectToDatabase() error {
	db, err := pgxpool.Connect(context.Background(), "")
	if err != nil {
		return err
	}

	DB = db
	log.Println("Connected to PostgreSQL.")

	return nil
}
