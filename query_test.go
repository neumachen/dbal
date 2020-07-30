package dbal

import (
	"database/sql"
	"testing"

	"github.com/c2fo/testify/require"
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
				d := &dbal{db: &mockDb{PrepareOk: false}}

				params := make(map[string]interface{})
				// act
				_, qryErr := d.Query(badSQLStmnt, params)

				// assertion
				require.Error(t, qryErr)
			},
		},
		{
			desc: "db.Prepare did not return an error",
			assertion: func(t *testing.T, desc string) {
				// arrangement
				db := testLoadDBAL(t)
				defer db.Close()
				params := map[string]interface{}{
					"first_name": "bearpig",
					"last_name":  "man",
					"address":    []byte(`{"test": "foo"}`),
				}

				// act
				_, qryErr := db.Query(insertCustomer, params)

				// assetion
				require.NoError(t, qryErr)

				// clean up
				_, qryErr = db.Query(deleteCustomer, params)
				require.NoError(t, qryErr)
			},
		},
		{
			desc: "no params passed",
			assertion: func(t *testing.T, desc string) {
				// arrangement
				db := testLoadDBAL(t)

				// act
				_, qryErr := db.Query(selectAllCustomers, nil)

				// assertion
				require.NoError(t, qryErr)
			},
		},
	}

	for _, test := range tests {
		test.assertion(t, test.desc)
	}
}
