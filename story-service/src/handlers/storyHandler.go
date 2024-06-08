package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"

	"cloud.google.com/go/vertexai/genai"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/nhat8002nguyen/story-of-media-be/story-service/src/models"
	"github.com/nhat8002nguyen/story-of-media-be/story-service/src/services"
)

func UploadData(c *gin.Context) {
	file, _ := c.FormFile("file")

	user_id := c.Query("user_id")
	session_id := c.Query("session_id")

	if file == nil || user_id == "" || session_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing values"})
		return
	}

	if err := services.SaveFileData(user_id, session_id, file); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := services.GenerateContentFromFile(c, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	content, err := parseContentResponse(resp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	story := models.Message{Content: content, Sender: "model"}
	err = services.SaveMessage(user_id, session_id, story)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"story": story})
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

	gemini := services.GenaiClient.GenerativeModel(services.ModelName)
	chat := gemini.StartChat()

	chat.History, err = services.LoadChatHistory(userID, sessionID)
	if err != nil {
		conn.WriteJSON(gin.H{"error": err.Error()})
		return
	}

	historyData, err := json.Marshal(chat.History)
	if err != nil {
		conn.WriteJSON(gin.H{"error": err.Error()})
		return
	}

	// Send the response back to the WebSocket client
	if err := conn.WriteMessage(websocket.TextMessage, historyData); err != nil {
		log.Println("write:", err)
		return
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
			userID, sessionID, models.Message{Sender: "user", Content: string(message)},
		); err != nil {
			log.Printf("error saving user message: %v", err)
		}
		if err := services.SaveMessage(
			userID, sessionID, models.Message{Sender: "model", Content: response},
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

func GetChatHistory(c *gin.Context) {
	user_id := c.Query("user_id")
	stories, err := services.GetStories(user_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"stories": stories})
}
