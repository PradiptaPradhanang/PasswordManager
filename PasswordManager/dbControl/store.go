package dbControl

import (
	"database/sql"
)

var DB *sql.DB

func Set(db *sql.DB) {
	DB = db
}

func Get() *sql.DB {
	return DB
}
