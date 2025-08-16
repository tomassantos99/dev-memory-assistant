package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "dev-memory-assistant.db")
	if err != nil {
		// handle error
	}
	defer db.Close()

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS clipboard (
    		id INTEGER PRIMARY KEY AUTOINCREMENT,
    		message TEXT,
    		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TRIGGER keep_last_100_rows
		AFTER INSERT ON clipboard
		WHEN (SELECT COUNT(*) FROM clipboard) > 100
		BEGIN
    		DELETE FROM clipboard
    		WHERE id IN (
        		SELECT id
        		FROM clipboard
        		ORDER BY id ASC
        		LIMIT 1
    		);
		END;
	`)

	if err != nil {
		fmt.Println("Error creating tables/triggers:", err)
	}
}
