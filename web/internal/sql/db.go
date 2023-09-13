package sql

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Scanner interface {
	Scan(dest ...interface{}) error
}

var dbpool *pgxpool.Pool = nil

func GetDatabasePool () *pgxpool.Pool {
    if dbpool != nil {
        return dbpool
    }

	DB_URL, isThere := os.LookupEnv("DB_URL")
	if !isThere || DB_URL == "" {
		log.Fatal("DB_URL env not set or has an empty value")
	}
	dbpool, err := pgxpool.New(context.Background(), DB_URL)
	if err != nil {
		log.Fatal("Could not create database connection pool: ", err)
	}

    return dbpool
}
