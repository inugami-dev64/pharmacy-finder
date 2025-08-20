package db

import (
	"fmt"
	"os"
	"pharmafinder"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

func EnsureMigrationsAreUpToDate(db *sqlx.DB) {
	goose.SetBaseFS(pharmafinder.MigrationsFS)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(db.DB, "db/migrations"); err != nil {
		panic(err)
	}
}

func ConnectToDB() *sqlx.DB {
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	db, err := sqlx.Connect("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	))

	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %v", err))
	}

	return db
}
