package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, data interface{}, msg string) {
	if msg == "" {
		msg = "success"
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": msg,
		"data":    data,
	})
}

func Error(c *gin.Context, errCode int, msg string) {
	c.JSON(errCode, gin.H{
		"status":  errCode,
		"message": msg,
		"data":    nil,
	})
}
