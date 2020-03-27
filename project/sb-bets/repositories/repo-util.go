package repo

import (
	"database/sql"
)

// Scanner comment
type Scanner interface {
	Scan(dest ...interface{}) error
}

// ScanRow scans one row and fills cells
func ScanRow(row Scanner, cells ...interface{}) error {
	err := row.Scan(cells...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	return nil
}
