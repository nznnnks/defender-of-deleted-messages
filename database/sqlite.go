package database

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./messages.db")
	if err != nil {
		return nil, err
	}

	query := `
	CREATE TABLE IF NOT EXISTS messages (
		chat_id INTEGER,
		message_id INTEGER,
		username TEXT,
		text TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

		PRIMARY KEY(chat_id, message_id)
	);
	`

	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}

	return db, nil
}
