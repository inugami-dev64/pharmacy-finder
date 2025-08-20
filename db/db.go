package db

import (
	"database/sql"
	"fmt"
	"os"
	"pharmafinder"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

// Query interface is a small abstraction that allows us
// to separate the logic, which creates the query and its
// execution
type Query[T any] interface {
	Query() (*T, error)
	QueryAll() ([]T, error)
	// Page method provides a capability to perform keyset paging
	// on any query.
	//
	// This is a prefered paging method due to performance reasonss
	Page(uniqueKey interface{}, key interface{}, length int, desc bool) ([]T, error)
}

type SQLXQuery[T any] struct {
	uniqueKey string
	key       string
	trx       *sqlx.DB
	q         string
	args      []interface{}
}

func (q *SQLXQuery[T]) Query() (*T, error) {
	var val T
	err := q.trx.Get(&val, q.q, q.args...)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &val, nil
}

func (q *SQLXQuery[T]) QueryAll() ([]T, error) {
	var vals []T
	err := q.trx.Select(&vals, q.q, q.args...)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return vals, nil
}

// This function needs testing on some actual data
// TODO: Write an integration test utilizing PostgreSQL testcontainer
// to test paging capability on some queries
func (q *SQLXQuery[T]) Page(uniqueKey interface{}, key interface{}, length int, desc bool) ([]T, error) {
	var outerQuery string
	args := q.args

	if !desc {
		outerQuery = fmt.Sprintf(`
		SELECT
		 	*
		FROM (
			%s
		) AS q
		WHERE
			(q."%s", q."%s") > (?, ?)
		ORDER BY
			q."%s",
			q."%s"
		LIMIT
			?
		`, q.q, q.key, q.uniqueKey, q.key, q.uniqueKey)
		args = append(args, key, uniqueKey, length)
	} else {
		outerQuery = fmt.Sprintf(`
		SELECT
		 	*
		FROM (
			%s
		) AS q
		WHERE
			(q."%s", q."%s") < (?, ?)
		ORDER BY
			q."%s" DESC,
			q."%s" DESC
		LIMIT
			?
		`, q.q, q.key, q.uniqueKey, q.key, q.uniqueKey)
		args = append(args, key, uniqueKey, length)
	}

	var data []T
	err := q.trx.Select(&data, outerQuery, args)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return data, nil
}

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
