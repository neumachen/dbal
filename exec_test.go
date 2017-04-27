package dbal

import (
	"database/sql"
	"testing"

	"github.com/magicalbanana/dbal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestExec(t *testing.T) {
	t.Parallel()
	tests := []struct {
		desc      string
		assertion func(*testing.T, string)
	}{
		{
			// db.Prepare errors out because the SQL file is empty
			desc: "db.Prepare returned an error",
			assertion: func(t *testing.T, desc string) {
				// arrangement
				d := &dbal{db: &mocks.Db{PrepareOk: false}}

				// act
				_, err := d.Exec(badSQLStmnt, nil)

				// assertion
				assert.Error(t, err)
			},
		},
		{
			desc: "db.Prepare did not return an error",
			assertion: func(t *testing.T, desc string) {
				db, openErr := sql.Open("postgres", "dbname=dbal_test sslmode=disable")
				assert.NoError(t, openErr)
				params := map[string]interface{}{
					"first_name": "bearpig",
					"last_name":  "man",
					"address":    []byte(`{ "test": "foo" }`),
				}
				d := &dbal{db: db}

				// act
				_, qryErr := d.Exec(insertCustomer, params)

				// assertion
				assert.NoError(t, qryErr)

				// clean up
				_, qryErr = d.Query(deleteCustomer, params)
				assert.NoError(t, qryErr)
			},
		},
	}

	for _, test := range tests {
		test.assertion(t, test.desc)
	}
}
