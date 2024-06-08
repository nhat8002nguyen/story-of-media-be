package services

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"cloud.google.com/go/vertexai/genai"
	"github.com/nhat8002nguyen/story-of-media-be/story-service/src/models"
	"google.golang.org/api/option"
)

const (
	projectID = "analyzing-media-files-web-app"
	location  = "asia-southeast1"
	ModelName = "gemini-1.5-flash-001"
)

var GenaiClient *genai.Client

func CreateGenAIClient() error {
	var ctx = context.Background()
	// Ensure the environment variable is set
	credsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credsFile == "" {
		log.Fatalf("GOOGLE_APPLICATION_CREDENTIALS environment variable not set")
	}
	var err error
	GenaiClient, err = genai.NewClient(ctx, projectID, location, option.WithCredentialsFile(credsFile))
	return err
}

func GenerateContentFromFile(
	c context.Context,
	file *multipart.FileHeader,
) (*genai.GenerateContentResponse, error) {
	f, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("could not open file")
	}
	defer f.Close()

	gemini := GenaiClient.GenerativeModel(ModelName)
	gemini.SetTemperature(1)

	fileBytes, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("unable to read file")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))

	prompt := "Generate a details story to describe this file"
	switch ext {
	case ".jpg":
		return gemini.GenerateContent(c, genai.Text(prompt), genai.ImageData("jpg", fileBytes))
	case ".jpeg":
		return gemini.GenerateContent(c, genai.Text(prompt), genai.ImageData("jpeg", fileBytes))
	case ".png":
		return gemini.GenerateContent(c, genai.Text(prompt), genai.ImageData("png", fileBytes))
	case ".pdf":
		return gemini.GenerateContent(c, genai.Text(prompt), genai.Blob{MIMEType: "application/pdf", Data: fileBytes})
	case ".txt":
		return gemini.GenerateContent(c, genai.Text(prompt), genai.Blob{MIMEType: "txt/plain", Data: fileBytes})
	default:
		return nil, fmt.Errorf("unknown or unsupported file format")
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

func GetStories(userID string) ([]string, error) {
	stmt := "SELECT session_id FROM chat_sessions WHERE user_id = $1"

	rows, err := models.Db.Query(stmt, userID)
	if err != nil {
		return nil, err
	}
	sessionIDs := []string{}
	sessionID := ""
	for rows.Next() {
		if err = rows.Scan(&sessionID); err != nil {
			return nil, err
		}
		sessionIDs = append(sessionIDs, sessionID)
	}
	return sessionIDs, nil
}
