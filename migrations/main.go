package main

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/pressly/goose/v3"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

//go:embed sql/*.sql
var embedMigrations embed.FS

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:\n./migrate <cmd> [migration_name | migration_version]")
		fmt.Println("Allowed cmds:\n- up\n- down\n- down_to <migration_version>")
		fmt.Println("- create <migration_name>\n- reset")
		os.Exit(1)
	}

	cmd := os.Args[1]

	if cmd == "create" && len(os.Args) < 3 {
		log.Fatal("When using create, must provide migration_name")
	}

	if cmd == "down_to" && len(os.Args) < 3 {
		log.Fatal("When using down_to, must provide migration_version")
	}

	dbUrl := os.Getenv("TURSO_DATABASE_URL")
	authToken := os.Getenv("TURSO_AUTH_TOKEN")

	if dbUrl == "" {
		log.Fatal("TURSO_DATABASE_URL env not set or has an empty value")
	}

	if authToken == "" {
		log.Fatal("TURSO_AUTH_TOKEN env not set or has an empty value")
	}

	sqlURL := fmt.Sprintf("%s?authToken=%s", dbUrl, authToken)
	db, err := sql.Open("libsql", sqlURL)
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("sqlite"); err != nil {
		log.Fatal(err)
	}

	migrationsDir := "sql"

	switch cmd {
	case "up":
		err = goose.Up(db, migrationsDir)
	case "down":
		err = goose.Down(db, migrationsDir)
	case "down_to":
		version, stdinErr := strconv.Atoi(os.Args[2])
		if stdinErr != nil {
			err = stdinErr
			break
		}
		err = goose.DownTo(db, migrationsDir, int64(version))
	case "create":
		fileName := os.Args[2]
		err = goose.Create(db, migrationsDir, fileName, "sql")
	case "reset":
		err = goose.Reset(db, migrationsDir)
	}

	if err != nil {
		log.Fatal(err)
	}
}
