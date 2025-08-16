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

func GetLastClipboardMessages(limit int) ([]string, error) {
	db, err := sql.Open("sqlite", "dev-memory-assistant.db")
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

