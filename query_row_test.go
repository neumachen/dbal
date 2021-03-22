package dbal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueryRow(t *testing.T) {
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
				d := &dbal{db: &mockDb{PrepareOk: false}}

				params := make(map[string]interface{})
				// act
				_, err := d.QueryRow(badSQLStmnt, params)

				// assertion
				require.Error(t, err)
			},
		},
		{
			desc: "db.Prepare did not return an error",
			assertion: func(t *testing.T, desc string) {
				// arrangement
				// arrangement
				db := testLoadDBAL(t)
				defer db.Close()
				params := map[string]interface{}{
					"first_name": "bearpig",
					"last_name":  "man",
					"address":    []byte(`{"test": "foo"}`),
				}

				// act
				_, qryErr := db.QueryRow(insertCustomer, params)

				// assertion
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
				defer db.Close()

				// act
				_, qryErr := db.QueryRow(selectAllCustomers, nil)

				// assertion
				require.NoError(t, qryErr)
			},
		},
	}

	for _, test := range tests {
		test.assertion(t, test.desc)
	}
}
