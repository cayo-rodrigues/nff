package database

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
	fredis "github.com/gofiber/storage/redis/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

var instance *Database

type Database struct {
	PG           *pgxpool.Pool
	Redis        *Redis
	SessionStore *session.Store
}

func (db *Database) Close() {
	if db.PG != nil {
		db.PG.Close()
	}
	if db.Redis.Client != nil {
		db.Redis.Close()
	}
	if db.SessionStore != nil {
		db.SessionStore.Storage.Close()
	}
}

// Estabilishes and returns a new DB connection
func NewDatabase() (*Database, error) {
	if instance != nil {
		instance.Close()
	}
	instance = new(Database)
	instance.Redis = new(Redis)

	err := initPG()
	if err != nil {
		return nil, err
	}
	err = initRedis()
	if err != nil {
		return nil, err
	}
	err = initSessionStore()
	if err != nil {
		return nil, err
	}

	return instance, nil
}

// Should be called only after NewDatabase is called, otherwise returns nil
func GetDB() *Database {
	return instance
}

func initPG() error {
	fmt.Println("Initializing PG connection...")

	if instance.PG != nil {
		fmt.Println("Reusing existing instance.PG connection")
		return nil
	}

	DB_URL := os.Getenv("DB_URL")
	if DB_URL == "" {
		return errors.New("DB_URL env missing or empty")
	}

	dbpool, err := pgxpool.New(context.Background(), DB_URL)
	if err != nil {
		return err
	}

	err = dbpool.Ping(context.Background())
	if err != nil {
		return errors.New(fmt.Sprintf("PG connection is not OK, ping failed: %v", err))
	}

	instance.PG = dbpool

	fmt.Println("New instance.PG connection OK")
	return nil
}

func initRedis() error {
	fmt.Println("Initiating Redis connection...")

	if instance.Redis.Client != nil {
		fmt.Println("Reusing existing instance.Redis connection")
		return nil
	}

	REDIS_URL := os.Getenv("REDIS_URL")
	if REDIS_URL == "" {
		return errors.New("REDIS_URL env missing or empty")
	}
	redisOpts, err := redis.ParseURL(REDIS_URL)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not parse redis url: %v", err))
	}
	rdb := redis.NewClient(redisOpts)

	err = rdb.Ping(context.Background()).Err()
	if err != nil {
		return errors.New(fmt.Sprintf("Redis connection is not OK, ping failed: %v", err))
	}

	instance.Redis.Client = rdb

	fmt.Println("New instance.Redis connection OK")
	return nil
}

func initSessionStore() error {
	fmt.Println("Initiating SessionStore connection...")

	if instance.SessionStore != nil {
		fmt.Println("Reusing existing instance.SessionStore connection")
		return nil
	}

	instance.SessionStore = session.New(session.Config{
		Storage: fredis.New(fredis.Config{
			URL: getRedisURL(),
		}),
		CookieHTTPOnly: true,
		Expiration:     7 * (24 * time.Hour),
	})

	fmt.Println("New instance.SessionStore connection OK")

	return nil
}

func getRedisURL() string {
	if instance.Redis == nil {
		return ""
	}

	redisOpts := instance.Redis.Options()
	return fmt.Sprintf("redis://%s/%d", redisOpts.Addr, redisOpts.DB)
}
