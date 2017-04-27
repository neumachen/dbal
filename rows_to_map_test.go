package dbal

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRowsToMap(t *testing.T) {
	t.Parallel()
	tests := []struct {
		desc      string
		assertion func(t *testing.T, desc string)
	}{
		{
			desc: "maps row columns to a map[string]interface{}",
			assertion: func(t *testing.T, desc string) {
				// arrangement
				db, openErr := sql.Open("postgres", "dbname=dbal_test sslmode=disable")
				assert.NoError(t, openErr)
				params := map[string]interface{}{
					"first_name": "piggly",
					"last_name":  "wiggly",
					"address":    []byte(`{"test": "pigglywiggly"}`),
				}
				d := &dbal{db: db}
				_, qryErr := d.Query(insertCustomer, params)
				assert.NoError(t, qryErr, "insert customer")
				rows, qryErr := d.Query(selectCustomer, params)
				assert.NoError(t, qryErr)

				// act
				m := RowsToMap(rows)

				// assertion
				assert.Equal(t, m[0]["first_name"], "piggly", desc)
				assert.Equal(t, m[0]["last_name"], "wiggly", desc)

				// clean up
				_, qryErr = d.Query(deleteCustomer, params)
				assert.NoError(t, qryErr, "delete customer")
			},
		},
	}

	for _, test := range tests {
		test.assertion(t, test.desc)
	}
}
