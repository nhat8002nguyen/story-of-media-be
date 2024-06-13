package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/nhat8002nguyen/story-of-media-be/story-service/src/middlewares"
	"github.com/nhat8002nguyen/story-of-media-be/story-service/src/models"
	"github.com/nhat8002nguyen/story-of-media-be/story-service/src/router"
	"github.com/nhat8002nguyen/story-of-media-be/story-service/src/services"

	"github.com/gin-gonic/gin"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file.")
	}
	// scripts.SeedData()
	models.InitDB()
	err = services.CreateGenAIClient()
	if err != nil {
		log.Fatalf("can not create genai client: %v", err)
	}
}

func main() {
	r := gin.Default()

	r.Use(middlewares.AuthMiddleware)

	router.SetupRouter(r)
	r.Run(":8080")
}
