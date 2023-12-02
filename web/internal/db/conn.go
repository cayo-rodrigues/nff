package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2/middleware/session"
	fredis "github.com/gofiber/storage/redis/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Scanner interface {
	Scan(dest ...interface{}) error
}

var PG *pgxpool.Pool
var Redis *redis.Client
var SessionStore *session.Store

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
	redisOpts, err := redis.ParseURL(REDIS_URL)
	if err != nil {
		log.Fatal("Could not create redis db connection: ", err)
	}
	rdb := redis.NewClient(redisOpts)

	err = rdb.Ping(context.Background()).Err()
	if err != nil {
		log.Fatal("Redis db connection is not OK, ping failed: ", err)
	}

	Redis = rdb
	return Redis
}

func GetSessionStore() *session.Store {
	if SessionStore != nil {
		return SessionStore
	}

	SessionStore = session.New(session.Config{
		Storage: fredis.New(fredis.Config{
			URL: getRedisURL(),
		}),
	})

	return SessionStore
}

func getRedisURL() string {
	if Redis == nil {
		return ""
	}

	redisOpts := Redis.Options()
	return fmt.Sprintf("redis://%s/%d", redisOpts.Addr, redisOpts.DB)
}
