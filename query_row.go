package dbal

import (
	"database/sql"

	"github.com/magicalbanana/dbal/sqltmpl"
)

// QueryRow uses the given sqlFile and params to execute sql.QueryRow. Before
// it gets executed the sqlFile is first retrieved from the virtual file
// system that was supplied and then passed to the sqltmpl to be parsed so
// that the named parameters are then replaced with positional parameters that
// are dependent on the driver that was used when initializing the DAL.
func (d *dbal) QueryRow(sqlStmnt string, params QueryParams) (*sql.Row, error) {
	// if params is not nil, used the parser
	if params != nil {
		tmpl := sqltmpl.NewParser(sqlStmnt)
		tmpl.SetValuesFromMap(params)

		stmt, err := d.db.Prepare(tmpl.GetParsedQuery())
		if err != nil {
			return nil, err
		}

		return stmt.QueryRow(tmpl.GetParsedParameters()...), nil
	}

	return d.db.QueryRow(sqlStmnt), nil
}
