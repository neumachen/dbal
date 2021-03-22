package dbal

import (
	"testing"

	"github.com/stretchr/testify/require"
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
				db := testLoadDBAL(t)
				defer db.Close()
				params := map[string]interface{}{
					"first_name": "piggly",
					"last_name":  "wiggly",
					"address":    []byte(`{"test": "pigglywiggly"}`),
				}

				_, qryErr := db.Query(insertCustomer, params)
				require.NoError(t, qryErr, "insert customer")
				rows, qryErr := db.Query(selectCustomer, params)
				require.NoError(t, qryErr)

				// act
				m, _ := RowsToMap(rows)

				// assertion
				require.Equal(t, m[0]["first_name"], "piggly", desc)
				require.Equal(t, m[0]["last_name"], "wiggly", desc)

				// clean up
				_, qryErr = db.Query(deleteCustomer, params)
				require.NoError(t, qryErr, "delete customer")
			},
		},
	}

	for _, test := range tests {
		test.assertion(t, test.desc)
	}
}
