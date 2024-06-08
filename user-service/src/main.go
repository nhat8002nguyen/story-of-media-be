package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/nhat8002nguyen/story-of-media-be/story-service/src/models"
	"github.com/nhat8002nguyen/story-of-media-be/story-service/src/router"

	"github.com/gin-gonic/gin"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file.")
	} else {
		log.Println("Loaded .env")
	}

	// scripts.SeedData()
	models.InitDB()

	r := gin.Default()
	router.SetupRouter(r)
	r.Run(":8081")
}
