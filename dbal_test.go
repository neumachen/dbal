package dbal

import (
	"database/sql"
	"log"
	"testing"

	"github.com/c2fo/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()
	db, openErr := sql.Open("postgres", testPgDbCreds(t).String())
	require.NoError(t, openErr)

	dal := New(db)
	require.NotNil(t, dal)
	defer db.Close()
}

func TestOpen(t *testing.T) {
	t.Parallel()
	tests := []struct {
		desc      string
		driver    string
		dbcreds   func() string
		assertion func(t *testing.T, e error, m string)
	}{
		{
			// db connection fails because of the invalid driver being passed
			// to sqlx.Open
			desc:   "database fails to open",
			driver: "manbearpig",
			dbcreds: func() string {
				return "dbname=piggly_wiggly"
			},
			assertion: func(t *testing.T, e error, m string) {
				require.Error(t, e, m)
			},
		},
		{
			desc:   "database loads successfully",
			driver: "postgres",
			dbcreds: func() string {
				return "dbname=piggly_wiggly"
			},
			assertion: func(t *testing.T, e error, m string) {
				require.NoError(t, e, m)
			},
		},
	}

	for _, test := range tests {
		db, err := Open(test.driver, test.dbcreds())
		defer func() {
			if db != nil {
				db.Close()
			}
		}()
		test.assertion(t, err, test.desc)
	}
}

func TestPingDataBase(t *testing.T) {
	t.Parallel()
	tests := []struct {
		desc      string
		driver    string
		db        func() DB
		assertion func(t *testing.T, e error, m string)
	}{
		{
			// DB fails to ping because the database does not exist
			desc:   "database ping reaches maximum ping time",
			driver: "postgres",
			db: func() DB {
				db, _ := Open("postgres", "dbname=piggly_wiggly")
				return db
			},
			assertion: func(t *testing.T, e error, m string) {
				require.Error(t, e, m)
			},
		},
		{
			desc:   "database loads successfully",
			driver: "postgres",
			db: func() DB {
				db, _ := Open("postgres", testPgDbCreds(t).String())
				return db
			},
			assertion: func(t *testing.T, e error, m string) {
				require.NoError(t, e, m)
			},
		},
	}

	for _, test := range tests {
		lgr := func(msg string) {
			log.Println(msg)
		}
		pingErr := PingDatabase(test.db(), 2, lgr)
		defer func() {
			if db := test.db(); db != nil {
				db.Close()
			}
		}()
		test.assertion(t, pingErr, test.desc)
	}
}
