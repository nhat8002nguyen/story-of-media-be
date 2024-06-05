package router

import (
	"github.com/nhat8002nguyen/story-of-media-be/src/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.POST("/upload", handlers.UploadData)
		api.GET("/story/:id", handlers.GetStory)
		api.POST("/story/:id/question", handlers.AskQuestion)
		api.GET("/stories", handlers.GetAllStories)
		api.POST("/connect-stories", handlers.ConnectStories)
	}
}
