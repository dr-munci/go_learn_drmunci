package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProgress(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Henüz implemente edilmedi"})
	return
}
