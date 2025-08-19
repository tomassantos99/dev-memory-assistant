package storage

import (
	"database/sql"

	"github.com/tomassantos99/dev-memory-assistant/paste/pkg"
	_ "modernc.org/sqlite"
)

func SaveClipboardMessage(message string) error {
	db, err := openDB()
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

func GetLastClipboardMessages(limit int) ([]string, error) {
	db, err := openDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(
		`SELECT message FROM clipboard ORDER BY id DESC LIMIT ?`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []string
	for rows.Next() {
		var message string
		if err := rows.Scan(&message); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func openDB() (*sql.DB, error) {
	dbPath, err := pkg.GetPathRelativeToExe("dev-memory-assistant.db")
	if err != nil {
		panic(err)
	}

	return sql.Open("sqlite", dbPath)
}
