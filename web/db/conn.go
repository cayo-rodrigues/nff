package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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
	fmt.Println("Initiating Database connection...")

	if PG != nil {
		fmt.Println("Reusing existing Database connection")
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

	fmt.Println("New Database connection OK")

	PG = dbpool
	return PG
}

func GetRedisConn() *redis.Client {
	fmt.Println("Initiating Redis connection...")

	if Redis != nil {
		fmt.Println("Reusing existing Redis connection")
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

	fmt.Println("New Redis connection OK")

	Redis = rdb
	return Redis
}

func GetSessionStore() *session.Store {
	fmt.Println("Initiating SessionStore connection...")

	if SessionStore != nil {
		fmt.Println("Reusing existing SessionStore connection")
		return SessionStore
	}

	SessionStore = session.New(session.Config{
		Storage: fredis.New(fredis.Config{
			URL: getRedisURL(),
		}),
		CookieHTTPOnly: true,
		Expiration:     7 * 24 * time.Hour,
	})

	fmt.Println("New SessionStore connection OK")

	return SessionStore
}

func getRedisURL() string {
	if Redis == nil {
		return ""
	}

	redisOpts := Redis.Options()
	return fmt.Sprintf("redis://%s/%d", redisOpts.Addr, redisOpts.DB)
}
