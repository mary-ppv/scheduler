package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Init() error {
	dbFile := os.Getenv("TODO_DBFILE")

	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	var err error
	DB, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		fmt.Printf("database file %s does not exist, will be created\n", dbFile)
	} else if err != nil {
		return fmt.Errorf("error checking database file: %w", err)
	}

	createTable(DB)

	return nil
}

func createTable(db *sql.DB) {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS scheduler (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
			date VARCHAR(8) NOT NULL DEFAULT '',
            title TEXT NOT NULL DEFAULT '',
            comment TEXT NOT NULL DEFAULT '',
            repeat VARCHAR(128) NOT NULL DEFAULT '' 
        )`,
		`CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date)`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			fmt.Printf("failed to execute query: %w", err)
		}
	}
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
