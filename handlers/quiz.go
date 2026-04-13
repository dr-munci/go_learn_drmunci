package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetQuiz(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Henüz implemente edilmedi"})
	return
}

func CreateQuiz(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Henüz implemente edilmedi"})
	return
}

func SubmitQuiz(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Henüz implemente edilmedi"})
	return
}
