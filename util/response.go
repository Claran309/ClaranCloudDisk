package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 通用API响应结构
// @Description 通用API响应格式
type Response struct {
	// 状态码
	Status int `json:"status" example:"200"`
	// 消息
	Message string `json:"message" example:"success"`
	// 数据
	Data interface{} `json:"data"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}, msg string) {
	if msg == "" {
		msg = "success"
	}
	c.JSON(http.StatusOK, Response{
		Status:  http.StatusOK,
		Message: msg,
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, errCode int, msg string) {
	c.JSON(errCode, Response{
		Status:  errCode,
		Message: msg,
		Data:    nil,
	})
}
