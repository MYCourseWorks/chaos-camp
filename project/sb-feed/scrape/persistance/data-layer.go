package persistance

import "database/sql"

// Database comment
type Database interface {
	Init(connectionString string) error
	Connection() *sql.DB
	GetDateTimeFormat() string
	Close() error
}

// DataLayer comment
type DataLayer struct {
	db Database
}

// NewDataLayer comment
func NewDataLayer(db Database) *DataLayer {
	return &DataLayer{db: db}
}
