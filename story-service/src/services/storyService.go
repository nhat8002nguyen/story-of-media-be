package services

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"path/filepath"

	"cloud.google.com/go/vertexai/genai"
	"github.com/nhat8002nguyen/story-of-media-be/story-service/src/models"
)

func SaveStory(story models.Story) (models.Story, error) {
	sqlStatement := `
        INSERT INTO stories (id, content)
        VALUES ($1, $2)
        RETURNING id`

	id := ""
	err := models.Db.QueryRow(sqlStatement, story.ID, story.Content).Scan(&id)
	if err != nil {
		return story, err
	}
	story.ID = id
	return story, nil
}

func GetStoryByID(id string) (models.Story, error) {
	sqlStatement := `SELECT id, content FROM stories WHERE id=$1;`
	var story models.Story
	row := models.Db.QueryRow(sqlStatement, id)
	switch err := row.Scan(&story.ID, &story.Content); err {
	case sql.ErrNoRows:
		return story, fmt.Errorf("no records found for story id %s", id)
	case nil:
		return story, nil
	default:
		return story, err
	}
}

// SaveMessage save message to PostgreSQL database
func SaveMessage(userID, sessionID string, message models.Message) error {
	stmt := "INSERT INTO chat_sessions(user_id, session_id, message, sender) VALUES ($1, $2, $3, $4)"
	_, err := models.Db.Exec(stmt, userID, sessionID, message.Content, message.Sender)
	return err
}

// LoadMessage load messages from PostgreSQL database
func LoadMessages(sessionID string) ([]models.Message, error) {
	stmt := "SELECT message, sender FROM chat_sessions WHERE session_id = $1 ORDER BY timestamp"
	rows, err := models.Db.Query(stmt, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var message models.Message
		if err := rows.Scan(&message.Content, &message.Sender); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

// LoadChatHistory loads chat history from PostgreSQL database
func LoadChatHistory(userID, sessionID string) ([]*genai.Content, error) {
	contentType := ""
	fileData := []byte{}
	stmt := "SELECT file_data, content_type FROM session_files WHERE user_id=$1 AND session_id=$2"
	if err := models.Db.QueryRow(stmt, userID, sessionID).Scan(&fileData, &contentType); err != nil {
		switch err {
		case sql.ErrNoRows:
			log.Printf("there is no file of user %s, and session %s", userID, sessionID)
		default:
			return nil, err
		}
	}

	stmt = "SELECT sender, message FROM chat_sessions WHERE user_id=$1 AND session_id=$2 ORDER BY timestamp"
	rows, err := models.Db.Query(stmt, userID, sessionID)
	if err != nil {
		return nil, err
	}

	// add media data first
	contents := []*genai.Content{
		{
			Role:  "user",
			Parts: []genai.Part{genai.Blob{MIMEType: contentType, Data: fileData}},
		},
	}

	for rows.Next() {
		var message models.Message
		if err := rows.Scan(&message.Sender, &message.Content); err != nil {
			return nil, err
		}
		contents = append(contents, &genai.Content{
			Role:  message.Sender,
			Parts: []genai.Part{genai.Text(message.Content)},
		})
	}
	return contents, nil
}

func SaveFileData(userID, sessionID string, file *multipart.FileHeader) error {
	filename := filepath.Base(file.Filename)
	contentType := file.Header.Get("content-type")
	f, err := file.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	fileData, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO session_files(user_id, session_id, filename, content_type, file_data)
	VALUES ($1, $2, $3, $4, $5) RETURNING id
	`

	return models.Db.QueryRow(stmt, userID, sessionID, filename, contentType, fileData).Err()
}
