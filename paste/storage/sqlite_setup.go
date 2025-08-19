package storage

import (
	"database/sql"
	"fmt"
	"github.com/tomassantos99/dev-memory-assistant/paste/pkg"
	_ "modernc.org/sqlite"
)

func SetupDB() {
	dbPath, err := pkg.GetPathRelativeToExe("dev-memory-assistant.db")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS clipboard (
    		id INTEGER PRIMARY KEY AUTOINCREMENT,
    		message TEXT,
    		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TRIGGER IF NOT EXISTS keep_last_100_rows
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
