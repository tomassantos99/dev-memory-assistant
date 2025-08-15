package storage

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

func SaveClipboardMessage(message string) error {

	db, err := sql.Open("sqlite", "dev-memory-assistant.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(
		`INSERT INTO clipboard (message) VALUES (?)`,
		message,
	)
	if err != nil {
		return err
	}

	return nil
}

