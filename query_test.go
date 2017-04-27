package dbal

import (
	"database/sql"
	"testing"

	"github.com/magicalbanana/dbal/mocks"

	"github.com/stretchr/testify/assert"
)

var db *sql.DB

func TestQuery(t *testing.T) {
	tests := []struct {
		desc      string
		assertion func(*testing.T, string)
	}{
		{
			desc: "db.Prepare returned an error",
			assertion: func(t *testing.T, desc string) {
				// arrangement
				d := &dbal{db: &mocks.Db{PrepareOk: false}}

				params := make(map[string]interface{})
				// act
				_, qryErr := d.Query(badSQLStmnt, params)

				// assertion
				assert.Error(t, qryErr)
			},
		},
		{
			desc: "db.Prepare did not return an error",
			assertion: func(t *testing.T, desc string) {
				// arrangement
				db, openErr := sql.Open("postgres", "dbname=dbal_test sslmode=disable")
				assert.NoError(t, openErr)
				params := map[string]interface{}{
					"first_name": "bearpig",
					"last_name":  "man",
					"address":    []byte(`{"test": "foo"}`),
				}
				d := &dbal{db: db}

				// act
				_, qryErr := d.Query(insertCustomer, params)

				// assetion
				assert.NoError(t, qryErr)

				// clean up
				_, qryErr = d.Query(deleteCustomer, params)
				assert.NoError(t, qryErr)
			},
		},
		{
			desc: "no params passed",
			assertion: func(t *testing.T, desc string) {
				// arrangement
				db, openErr := sql.Open("postgres", "dbname=dbal_test sslmode=disable")
				assert.NoError(t, openErr)
				d := &dbal{db: db}

				// act
				_, qryErr := d.Query(selectAllCustomers, nil)

				// assertion
				assert.NoError(t, qryErr)
			},
		},
	}

	for _, test := range tests {
		test.assertion(t, test.desc)
	}
}
