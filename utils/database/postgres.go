package database

import (
	"api.ethscrow/utils"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

var DB *pgxpool.Pool

func ConnectToDatabase() error {
	db, err := pgxpool.Connect(context.Background(), utils.DATABASE_URL)
	if err != nil {
		return err
	}

	DB = db
	log.Println("Connected to PostgreSQL.")

	return nil
}
