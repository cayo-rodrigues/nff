package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Scanner interface {
	Scan(dest ...interface{}) error
}

var PG *pgxpool.Pool = nil
var Redis *redis.Client = nil

func GetDBPool() *pgxpool.Pool {
	if PG != nil {
		return PG
	}

	DB_URL, isThere := os.LookupEnv("DB_URL")
	if !isThere || DB_URL == "" {
		log.Fatal("DB_URL env not set or has an empty value")
	}
	dbpool, err := pgxpool.New(context.Background(), DB_URL)
	if err != nil {
		log.Fatal("Could not create database connection pool: ", err)
	}

	err = dbpool.Ping(context.Background())
	if err != nil {
		log.Fatal("Database connection is not OK, ping failed: ", err)
	}

	PG = dbpool
	return PG
}

func GetRedisConn() *redis.Client {
	if Redis != nil {
		return Redis
	}

	REDIS_URL, isThere := os.LookupEnv("REDIS_URL")
	if !isThere || REDIS_URL == "" {
		log.Fatal("REDIS_URL env not set or has an empty value")
	}
	redis_opts, err := redis.ParseURL(REDIS_URL)
	if err != nil {
		log.Fatal("Could not create redis db connection: ", err)
	}
	rdb := redis.NewClient(redis_opts)

	err = rdb.Ping(context.Background()).Err()
	if err != nil {
		log.Fatal("Redis db connection is not OK, ping failed: ", err)
	}

	Redis = rdb
	return Redis
}
