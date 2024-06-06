package handlers

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"cloud.google.com/go/vertexai/genai"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/nhat8002nguyen/story-of-media-be/story-service/src/models"
	"github.com/nhat8002nguyen/story-of-media-be/story-service/src/services"
	"google.golang.org/api/option"
)

const (
	projectID = "analyzing-media-files-web-app"
	location  = "asia-southeast1"
	modelName = "gemini-1.5-flash-001"
)

var genaiClient *genai.Client

func CreateGenAIClient() error {
	var ctx = context.Background()
	// Ensure the environment variable is set
	credsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credsFile == "" {
		log.Fatalf("GOOGLE_APPLICATION_CREDENTIALS environment variable not set")
	}
	var err error
	genaiClient, err = genai.NewClient(ctx, projectID, location, option.WithCredentialsFile(credsFile))
	return err
}

func UploadData(c *gin.Context) {
	file, _ := c.FormFile("file")
	filename := filepath.Base(file.Filename)

	user_id := c.Query("user_id")
	session_id := c.Query("session_id")

	// Saving the uploaded file to local disk (or temporary storage)
	uploadPath := fmt.Sprintf("/uploads/%s", filename)
	c.SaveUploadedFile(file, uploadPath)

	gemini := genaiClient.GenerativeModel(modelName)
	gemini.SetTemperature(1)
	resp, err := generateContentFromFile(c, gemini, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	content, err := parseContentResponse(resp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	story := models.Message{Content: content, Sender: "ai"}
	err = services.SaveMessage(user_id, session_id, story)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"story": story})
}

func generateContentFromFile(
	c *gin.Context,
	gemini *genai.GenerativeModel,
	file *multipart.FileHeader,
) (*genai.GenerateContentResponse, error) {
	f, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("could not open file")
	}
	defer f.Close()

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

// WebSocket Upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WsHandler is WebSocket handler function
func WsHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("upgrade: ", err)
		return
	}
	defer conn.Close()

	userID := c.Query("user_id")
	sessionID := c.Query("session_id")

	gemini := genaiClient.GenerativeModel(modelName)
	chat := gemini.StartChat()

	// Load previous messages from the database
	messages, err := services.LoadMessages(sessionID)
	if err != nil {
		log.Printf("error loading messages: %v", err)
		return
	}
	for _, msg := range messages {
		// Check if the message is from the client
		if msg.Sender == "client" {
			// Send historical client messages to the gemini chat to maintain context
			_, err := makeChatRequests(c, chat, msg.Content)
			if err != nil {
				log.Printf("error sending historical message: %v", err)
				return
			}
		}
	}

	for {
		// Read message from WebSocket
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		// Make chat request with received message
		response, err := makeChatRequests(c, chat, string(message))
		if err != nil {
			log.Printf("makeChatRequests error: %v", err)
			break
		}

		// Save the user message and the response to the database
		if err := services.SaveMessage(
			userID, sessionID, models.Message{Sender: "client", Content: string(message)},
		); err != nil {
			log.Printf("error saving user message: %v", err)
		}
		if err := services.SaveMessage(
			userID, sessionID, models.Message{Sender: "ai", Content: response},
		); err != nil {
			log.Printf("error saving response message: %v", err)
		}

		// Send the response back to the WebSocket client
		if err := conn.WriteMessage(websocket.TextMessage, []byte(response)); err != nil {
			log.Println("write:", err)
			break
		}
	}
}

// makeChatRequests send chat request to the Gemini model
func makeChatRequests(ctx context.Context, chat *genai.ChatSession, message string) (string, error) {
	r, err := chat.SendMessage(ctx, genai.Text(message))
	if err != nil {
		return "", err
	}
	return parseContentResponse(r)
}

func parseContentResponse(r *genai.GenerateContentResponse) (string, error) {
	part := r.Candidates[0].Content.Parts[0]
	value := reflect.ValueOf(part)
	if value.Kind() == reflect.String {
		return value.String(), nil
	} else {
		return "", fmt.Errorf("not found response text")
	}
}
