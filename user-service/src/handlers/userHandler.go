package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nhat8002nguyen/story-of-media-be/story-service/src/models"
	"github.com/nhat8002nguyen/story-of-media-be/story-service/src/services"
)

func GetUserByEmail(c *gin.Context) {
	email := c.Param("email")

	u, err := services.GetUser(email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": u})
}

func AddUser(c *gin.Context) {
	user := &models.User{}
	if err := c.ShouldBindJSON(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	user, err := services.AddUser(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "id": user.ID})
}

func LoginHanlder(c *gin.Context) {
	u := models.User{}
	err := c.ShouldBindJSON(&u)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	token, err := services.Authenticate(u)
	if err != nil {
		if customErr, ok := err.(*services.CustomError); ok {
			if customErr.Code == services.ERROR_NOT_FOUND {
				c.JSON(http.StatusUnauthorized, gin.H{"error": customErr.Message})
				return
			} else if customErr.Code == services.ERROR_UNAUTHORIZED {
				c.JSON(http.StatusUnauthorized, gin.H{"error": customErr.Message})
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": customErr.Message})
			return
		}
	}

	c.SetCookie("token", token, 60*60, "/", "story-of-media-ai.vercel.app", true, true)
}
