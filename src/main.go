package main

import (
	"github.com/nhat8002nguyen/story-of-media-be/src/models"
	"github.com/nhat8002nguyen/story-of-media-be/src/router"

	"github.com/gin-gonic/gin"
)

func main() {
	models.InitDB()
	r := gin.Default()
	router.SetupRouter(r)
	r.Run(":8080")
}
