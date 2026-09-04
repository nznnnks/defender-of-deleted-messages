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
	mediaType string,
	fileID string,
	caption string,
) error {
	_, err := db.Exec(
		`INSERT INTO messages(chat_id, message_id, username, text, mediaType, fileID, caption)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		chatID,
		messageID,
		username,
		text,
		mediaType,
		fileID,
		caption,
	)

	return err
}

func GetMessage(
	db *sql.DB,
	chatID int64,
	messageID int,
) (string, //username
	string, //text
	string, //mediaType
	string, //fileID
	string, //caption
	error,
) {
	var username string
	var text string
	var mediaType string
	var fileID string
	var caption string

	err := db.QueryRow(
		`SELECT username, text, mediaType, fileID, caption
		FROM messages
		WHERE chat_id = ?
		AND message_id = ?`,
		chatID,
		messageID,
	).Scan(&username,
		&text,
		&mediaType,
		&fileID,
		&caption,
	)

	if err != nil {
		return "", "", "", "", "", err
	}

	return username, text, mediaType, fileID, caption, nil
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
