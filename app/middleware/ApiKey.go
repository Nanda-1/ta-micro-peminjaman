package middleware

import (
	"log"
	"net/http"
	"os"
	"ta_microservice_peminjaman/app/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func ApiKey() gin.HandlerFunc {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ServiceAPIKey := os.Getenv("API.KEY")

	if ServiceAPIKey == "" {
		log.Fatal("Please set API.KEY environment variable")
	}

	return func(c *gin.Context) {
		res := models.JsonResponse{Success: true}
		apiKey := c.Request.Header.Get("API.KEY")

		var errorMsg string

		if apiKey == "" {
			errorMsg = "Missing Key"
			res.Success = false
			res.Error = &errorMsg
			c.JSON(http.StatusBadRequest, res)
			c.Abort()
			return
		} else if ServiceAPIKey != apiKey {
			errorMsg = "Invalid Key"
			res.Success = false
			res.Error = &errorMsg
			c.JSON(http.StatusUnauthorized, res)
			c.Abort()
			return
		}

		c.Next()
	}
}
