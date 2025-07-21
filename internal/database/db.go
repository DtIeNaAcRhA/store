package database

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB(filepath string) error {
	var err error
	DB, err = sql.Open("sqlite", filepath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	_, err = DB.Exec(`PRAGMA foreign_keys = ON`)
	if err != nil {
		return fmt.Errorf("failed to enable foreign keys: %w", err)
	}
	return nil

}
