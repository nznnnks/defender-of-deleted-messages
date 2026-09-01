package database

import (
	"database/sql"
)

func SaveMessage(
	db *sql.DB,
	chatID int64,
	messageID int,
	username string,
	text string,
) error {
	_, err := db.Exec(
		`INSERT INTO messages(chat_id, message_id, username, text)
		VALUES (?, ?, ?, ?)`,
		chatID,
		messageID,
		username,
		text,
	)

	return err
}

func GetMessage(
	db *sql.DB,
	chatID int64,
	messageID int,
) (string, string,  error) {
	var username string
	var text string

	err := db.QueryRow(
		`SELECT text
		FROM messages
		WHERE chat_id = ?
		AND message_id = ?`,
		chatID,
		messageID,
	).Scan(&username, &text)

	if err != nil {
		return "", "", err
	}

	return username, text, nil
}

func UpdateMessage(
	db *sql.DB,
	chatID int64,
	messageID int,
	text string,
) error {
	_, err := db.Exec(
		`UPDATE messages
		SET text = ?
		WHERE chat_id = ?
		AND message_id = ?`,
		text,
		chatID,
		messageID,
	)
	return err
}
