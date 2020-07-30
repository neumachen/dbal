package dbal

import (
	"database/sql"
	"fmt"
	"time"

	// postgres driver
	_ "github.com/lib/pq"
)

// DBAL (database access layer) is the type returned when initializing an new
// instance of this package. The functions listed are the functions that can
// be used for this package.
type DBAL interface {
	Exec(sqlFile string, params QueryParams) (sql.Result, error)
	Query(sqlFile string, params QueryParams) (*sql.Rows, error)
	QueryRow(sqlFile string, params QueryParams) (*sql.Row, error)
	Close() error
}

// DB is an interface modeled after the go's standard database/sql package.
type DB interface {
	Prepare(query string) (*sql.Stmt, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Ping() error
	Close() error
}

type dbal struct {
	db DB
}

// Close ...
func (d *dbal) Close() error {
	return d.db.Close()
}

// New returns a DBAL using the given DB and FileStore.
func New(db DB) DBAL {
	return &dbal{
		db: db,
	}
}

// Open opens a new database connection using the given driver and the
// dbCreds. If an error occurs when opening a connection an error is
// returned.
func Open(driver string, dbCreds string) (DB, error) {
	db, err := sql.Open(driver, dbCreds)
	if err != nil {
		return nil, fmt.Errorf("database connection failed")
	}
	return db, nil
}

// PingDatabase is a helper function to ping the database with backoff
// to ensure a connection can be established before we proceed with a
// database setup whilst logging each ping and returns an error if it fails
// using the given function to log.
func PingDatabase(db DB, dbPingTime int, lgr func(msg string)) error {
	for i := 0; i < dbPingTime; i++ {
		err := db.Ping()
		if err == nil {
			return err
		}
		if lgr != nil {
			lgr("Database ping failed. Retry in 2s")
		}
		time.Sleep(time.Second)
	}
	return fmt.Errorf("failed to ping the database after %v seconds", dbPingTime)
}
