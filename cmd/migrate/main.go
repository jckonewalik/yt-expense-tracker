package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", "ytexpensetracker", "admin@123", "localhost", "ytexpensetracker"))
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("db not responding: %v", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("error getting db instance: %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"postgres", driver)

	if err != nil {
		log.Fatalf("error getting migrate db instance: %v", err)
	}

	if err := m.Up(); err != nil {
		log.Fatalf("error running migrations: %v", err)
	}
}
