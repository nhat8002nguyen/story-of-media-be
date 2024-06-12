package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nhat8002nguyen/story-of-media-be/story-service/src/models"
	"github.com/nhat8002nguyen/story-of-media-be/story-service/src/router"
	"github.com/nhat8002nguyen/story-of-media-be/story-service/src/services"

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
	services.SecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

	r := gin.Default()
	router.SetupRouter(r)
	r.Run(":8081")
}
