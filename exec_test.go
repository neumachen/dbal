package dbal

import (
	"testing"

	"github.com/stretchr/testify/require"
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
				d := &dbal{db: &mockDb{PrepareOk: false}}

				// act
				_, err := d.Exec(badSQLStmnt, nil)

				// assertion
				require.Error(t, err, desc)
			},
		},
		{
			desc: "db.Prepare did not return an error",
			assertion: func(t *testing.T, desc string) {
				db := testLoadDBAL(t)
				defer db.Close()
				params := map[string]interface{}{
					"first_name": "bearpig",
					"last_name":  "man",
					"address":    []byte(`{ "test": "foo" }`),
				}

				// act
				_, qryErr := db.Exec(insertCustomer, params)

				// assertion
				require.NoError(t, qryErr, desc)

				// clean up
				_, qryErr = db.Query(deleteCustomer, params)
				require.NoError(t, qryErr, desc)
			},
		},
	}

	for _, test := range tests {
		test.assertion(t, test.desc)
	}
}
