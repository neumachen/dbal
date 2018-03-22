package dbal

import (
	"database/sql"

	"github.com/magicalbanana/dbal/sqltmpl"
)

func (d *dbal) Exec(sqlStmnt string, params QueryParams) (sql.Result, error) {
	tmpl := sqltmpl.NewParser(sqlStmnt)
	tmpl.SetValuesFromMap(params)

	stmt, err := d.db.Prepare(tmpl.GetParsedQuery())
	if err != nil {
		return nil, err
	}
	return stmt.Exec(tmpl.GetParsedParameters()...)
}
