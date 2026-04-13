package handlers

import (
	"github.com/gin-gonic/gin"
)

func ClassroomWS(c *gin.Context) {
	c.JSON(501, gin.H{"error": "WebSocket henüz implemente edilmedi"})
	return
}
