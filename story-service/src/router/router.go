package router

import (
	"github.com/nhat8002nguyen/story-of-media-be/story-service/src/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.POST("/upload", handlers.UploadData)
		api.GET("/story/ws", handlers.WsHandler)
		api.GET("/stories", handlers.GetChatHistory)
	}
}
