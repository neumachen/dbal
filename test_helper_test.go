package dbal

import (
	"net/url"
	"os"
	"testing"

	"github.com/c2fo/testify/require"
)

var badSQLStmnt = `SSSSSS`
var insertCustomer = `insert into customers(first_name, last_name, address) values($first_name, $last_name, CAST(NULLIF($address, '') as jsonb));`
var selectCustomer = `select first_name, last_name, address from customers where first_name = $first_name AND last_name = $last_name;`
var selectAllCustomers = `select first_name, last_name, address from customers;`
var deleteCustomer = `delete from customers where first_name = $first_name and last_name = $last_name;`

func testPgDbCreds(t *testing.T) *url.URL {
	v := os.Getenv("DBAL_PG_DB")
	if v == "" {
		v = "postgres://dbal:dbal@localhost:5432/dbal_development?sslmode=disable"
	}
	u, err := url.Parse(v)
	require.NoError(t, err)
	return u
}

func testLoadDBAL(t *testing.T) DBAL {
	db, err := Open("postgres", testPgDbCreds(t).String())
	require.NoError(t, err)
	return &dbal{
		db: db,
	}
}
