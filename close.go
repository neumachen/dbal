package dbal

// Close ...
func (d *dbal) Close() error {
	return d.db.Close()
}
