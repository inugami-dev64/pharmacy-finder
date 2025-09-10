package db

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"pharmafinder"
	"pharmafinder/utils"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"

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
	vals := []T{}
	err := q.trx.Select(&vals, q.q, q.args...)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return vals, nil
}

// Utility function, which extracts pager
// HTTP query parameters and returns them
//
// Pager variables are following:
//   - uk - specifying unique key value
//   - k - specifying key value
//   - l - specifying query set length
//   - desc - specifying whether the pager should work in descending order
func ExtractPagerQueryParameters(params url.Values) (string, string, int, bool) {
	uk := params.Get("uk")
	k := params.Get("k")
	lStr := params.Get("l")
	descStr := params.Get("desc")

	var err error
	var l int64
	var desc bool

	if l, err = strconv.ParseInt(lStr, 10, 64); err != nil {
		l = 50
	}

	if desc, err = strconv.ParseBool(descStr); err != nil {
		desc = false
	}

	return uk, k, int(l), desc
}

// This function needs testing on some actual data
// TODO: Write an integration test utilizing PostgreSQL testcontainer
// to test paging capability on some queries
func (q *SQLXQuery[T]) Page(uniqueKey interface{}, key interface{}, length int, desc bool) ([]T, error) {
	var outerQuery string
	args := q.args
	c := len(args) + 1

	if uniqueKey == nil || key == nil {
		if !desc {
			outerQuery = fmt.Sprintf(
				`SELECT
					*
				FROM (
					%s
				) AS q
				ORDER BY
					q."%s",
					q."%s"
				LIMIT
					$%d
				`, q.q, q.key, q.uniqueKey, c)
			args = append(args, length)
		} else {
			outerQuery = fmt.Sprintf(
				`SELECT
					*
				FROM (
					%s
				) AS q
				ORDER BY
					q."%s" DESC,
					q."%s" DESC
				LIMIT
					$%d
				`, q.q, q.key, q.uniqueKey, c)
			args = append(args, length)
		}
	} else {
		if !desc {
			outerQuery = fmt.Sprintf(
				`SELECT
					*
				FROM (
					%s
				) AS q
				WHERE
					(q."%s", q."%s") > ($%d, $%d)
				ORDER BY
					q."%s",
					q."%s"
				LIMIT
					$%d
				`, q.q, q.key, q.uniqueKey, c, c+1, q.key, q.uniqueKey, c+2)
			args = append(args, key, uniqueKey, length)
		} else {
			outerQuery = fmt.Sprintf(
				`SELECT
					*
				FROM (
					%s
				) AS q
				WHERE
					(q."%s", q."%s") < ($%d, $%d)
				ORDER BY
					q."%s" DESC,
					q."%s" DESC
				LIMIT
					$%d
				`, q.q, q.key, q.uniqueKey, c, c+1, q.key, q.uniqueKey, c+2)
			args = append(args, key, uniqueKey, length)
		}
	}

	data := []T{}
	err := q.trx.Select(&data, outerQuery, args...)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return data, nil
}

// Custom logger type for goose logging
type GooseLogger struct {
	Logger zerolog.Logger
}

func (l *GooseLogger) Fatalf(format string, v ...interface{}) {
	l.Logger.Error().Msgf(format, v...)
}

func (l *GooseLogger) Printf(format string, v ...interface{}) {
	l.Logger.Info().Msgf(format, v...)
}

func EnsureMigrationsAreUpToDate(db *sqlx.DB) {
	goose.SetBaseFS(pharmafinder.MigrationsFS)
	goose.SetLogger(&GooseLogger{Logger: utils.GetLogger("DB")})

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(db.DB, "db/migrations"); err != nil {
		panic(err)
	}
}

// Attempts to connect to the database and
// update SQL migrations
func ProvideDatabaseHandle() *sqlx.DB {
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=utc",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %v", err))
	}

	// prepare the logger
	loggerOptions := []sqldblogger.Option{
		sqldblogger.WithSQLQueryFieldname("query"),
		sqldblogger.WithDurationUnit(sqldblogger.DurationMillisecond),
		sqldblogger.WithDurationFieldname("duration"),
		sqldblogger.WithQueryerLevel(sqldblogger.LevelDebug),
		sqldblogger.WithExecerLevel(sqldblogger.LevelDebug),
		sqldblogger.WithPreparerLevel(sqldblogger.LevelDebug),
	}
	adapter := zerologadapter.New(utils.GetLogger("SQL"))
	db = sqldblogger.OpenDriver(dsn, db.Driver(), adapter, loggerOptions...)

	// pass it to sqlx
	sqlxDB := sqlx.NewDb(db, "postgres")
	EnsureMigrationsAreUpToDate(sqlxDB)
	return sqlxDB
}
