package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/greet", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, welcome to our API!",
		})
	})
	router.POST("/submit", func(c *gin.Context) {
		var request struct {
			Name string `json:"name" binding:"required"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{
			"message": "Hello, " + request.Name + "!",
		})
	})
	router.Run(":8080")
}
