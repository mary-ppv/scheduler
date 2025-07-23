package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

func Init() error {
	dbFile := os.Getenv("TODO_DBFILE")

	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		fmt.Println("Database file does not exist")
	} else if err != nil {
		return err
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		createTable(db)
	}

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
			fmt.Printf("failed to execute query: %v", err)
		}
	}
}
