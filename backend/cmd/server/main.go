package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"app":    "RTMS_Monitoring",
		})
	})

	if err := router.Run("127.0.0.1:8080"); err != nil {
		log.Fatal(err)
	}
}
