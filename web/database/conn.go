package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
	fredis "github.com/gofiber/storage/redis/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

var instance *Database

type Database struct {
	SQLite       *sql.DB
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

	err := initSQLite()
	if err != nil {
		return nil, err
	}

	err = initPG()
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

// Should be called only after NewDatabase is called, otherwise returns nil
func GetSQLite() *sql.DB {
	db := GetDB()
	if db != nil {
		return db.SQLite
	}
	return nil
}

// Should be called only after NewDatabase is called, otherwise returns nil
func GetPG() *pgxpool.Pool {
	db := GetDB()
	if db != nil {
		return db.PG
	}
	return nil
}

// Should be called only after NewDatabase is called, otherwise returns nil
func GetRedis() *Redis {
	db := GetDB()
	if db != nil {
		return db.Redis
	}
	return nil
}

// Should be called only after NewDatabase is called, otherwise returns nil
func GetSessionStore() *session.Store {
	db := GetDB()
	if db != nil {
		return db.SessionStore
	}
	return nil
}

func initSQLite() error {
	fmt.Println("Initializing SQLite connection...")
	if instance.SQLite != nil {
		fmt.Println("Reusing existing instance.SQLite connection")
		return nil
	}

	tursoDatabaseUrl := os.Getenv("TURSO_DATABASE_URL")
	if tursoDatabaseUrl == "" {
		return errors.New("TURSO_DATABASE_URL env missing or empty")
	}

	authToken := os.Getenv("TURSO_AUTH_TOKEN")
	if authToken == "" {
		return errors.New("TURSO_AUTH_TOKEN env missing or empty")
	}

	fullURL := fmt.Sprintf("%s?authToken=%s", tursoDatabaseUrl, authToken)

	sqliteDB, err := sql.Open("libsql", fullURL)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to open DB %s: %v\n", fullURL, err))
	}

	if err := sqliteDB.Ping(); err != nil {
		return errors.New(fmt.Sprintf("SQLite Database connection is not OK, ping failed:", err))
	}

	instance.SQLite = sqliteDB
	fmt.Println("New instance.SQLite connection OK")
	return nil
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
