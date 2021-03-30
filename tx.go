package dbal

// Tx ...
// func (d *dbal) Tx(fn func(tx *sql.Tx) (*sql.Rows, error)) (*sql.Rows, error) {
// 	var rollback bool
// 	tx := d.db.Begin()
// 	defer func() {
// 		if rollback {
// 			tx.Rollback()
// 		}
// 		tx.Commit()
// 	}()
// 	sqlRows, err := fn(tx)
// 	if err != nil {
//
// 	}
//
// 	return nil, nil
// }
