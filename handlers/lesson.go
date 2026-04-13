package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetLessons(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Henüz implemente edilmedi"})
	return
}

func CreateLesson(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Henüz implemente edilmedi"})
	return
}

func CompleteLesson(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Henüz implemente edilmedi"})
	return
}
