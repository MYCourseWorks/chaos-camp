package mysqlpers

import (
	"database/sql"

	// Mysql driver should be imported, only if we use this implementation of OddsDb
	_ "github.com/go-sql-driver/mysql"
)

// DatabaseMySQL comment
type DatabaseMySQL struct {
	connection *sql.DB
}

// Init comment
func (db *DatabaseMySQL) Init(connectionString string) error {
	newDb, err := sql.Open("mysql", connectionString)
	if err != nil {
		return err
	}

	*db = DatabaseMySQL{
		connection: newDb,
	}
	return nil
}

// Close commnet
func (db *DatabaseMySQL) Close() error {
	return db.connection.Close()
}

// Connection comment
func (db *DatabaseMySQL) Connection() *sql.DB {
	return db.connection
}

// GetDateTimeFormat comment
func (db *DatabaseMySQL) GetDateTimeFormat() string {
	return "2006-01-02 15:04:05"
}
