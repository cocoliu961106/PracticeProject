package service

import "github.com/gin-gonic/gin"

// 业务逻辑
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ping",
	})
}
