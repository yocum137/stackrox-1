package postgres

import (
	"github.com/jackc/pgx/v4"
)

type Rows struct {
	rowsScanned int
	pgx.Rows
}

func (r *Rows) Close() {
	// Eventually extended for metrics to report number of returned rows
	r.Rows.Close()
}

func (r *Rows) Scan(dest ...interface{}) error {
	err := r.Rows.Scan(dest...)
	if err != nil {
		return err
	}
	r.rowsScanned++
	return nil
}
