package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/nhat8002nguyen/story-of-media-be/story-service/src/models"
	"github.com/nhat8002nguyen/story-of-media-be/story-service/src/services"
)

func UploadData(c *gin.Context) {
	file, _ := c.FormFile("file")
	filename := filepath.Base(file.Filename)

	// Saving the uploaded file to local disk (or temporary storage)
	uploadPath := fmt.Sprintf("/uploads/%s", filename)
	c.SaveUploadedFile(file, uploadPath)

	// Invoke Gemini Vertex API for processing the file (this could be a separate service call)
	services.GenerateStory(file)

	story := models.Story{Content: "example content"}
	savedStory, err := services.SaveStory(story)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"story": savedStory})
}

func GetStory(c *gin.Context) {
	id := c.Param("id")

	story, err := services.GetStoryByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"story": story})
}

func AskQuestion(c *gin.Context) {
	// id := c.Param("id")
	var question struct {
		Text string `json:"text"`
	}
	c.ShouldBindJSON(&question)

	// Process question (invoke Gemini Vertex API or NLP model)
	// Retrieve or generate answer from the stored story

	c.JSON(http.StatusOK, gin.H{"answer": "Generated answer to the question."})
}

func GetAllStories(c *gin.Context) {
	// Query PostgreSQL for all stories

	c.JSON(http.StatusOK, gin.H{"stories": "List of stories"})
}

func ConnectStories(c *gin.Context) {
	var connection struct {
		StoryAID string `json:"story_a_id"`
		StoryBID string `json:"story_b_id"`
	}

	c.ShouldBindJSON(&connection)

	// Logic to connect stories in PostgreSQL

	c.JSON(http.StatusOK, gin.H{"message": "Stories connected successfully."})
}
