package router

import (
	"github.com/nhat8002nguyen/story-of-media-be/story-service/src/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/user/:email", handlers.GetUserByEmail)
		api.POST("/user", handlers.AddUser)
		api.POST("/login", handlers.LoginHanlder)
	}
}
