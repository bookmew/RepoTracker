package repository

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
)

func InitDB() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
		return nil, err
	}
	return conn, nil
}