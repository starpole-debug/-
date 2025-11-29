package response

import "github.com/gin-gonic/gin"

// Success standardises payload wrapper for API success responses.
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{"data": data})
}

// Created returns a 201 response for POST endpoints.
func Created(c *gin.Context, data interface{}) {
	c.JSON(201, gin.H{"data": data})
}

// Error centralises error serialisation for clients.
func Error(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"error": msg})
}
