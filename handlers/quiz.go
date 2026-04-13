package handlers

import (
	"fmt"
	"golearn/database"
	"golearn/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetQuiz(c *gin.Context) {
	var quiz models.Quiz
	if err := database.DB.Preload("Questions").Where("lesson_id = ?", c.Param("id")).First(&quiz).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiz bulunamadı"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": quiz})
	return
}

func CreateQuiz(c *gin.Context) {
	var quiz models.Quiz
	if err := c.ShouldBindJSON(&quiz); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var lesson models.Lesson
	if err := database.DB.First(&lesson, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ders bulunamadı"})
		return
	}

	quiz.LessonID = lesson.ID
	if err := database.DB.Create(&quiz).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Quiz oluşturulamadı"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Quiz oluşturuldu", "data": quiz})
	return
}

type SubmitAnswer struct {
	Answers map[string]string `json:"answers"`
}

func SubmitQuiz(c *gin.Context) {
	var input SubmitAnswer
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var quiz models.Quiz
	if err := database.DB.Preload("Questions").First(&quiz, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiz bulunamadı"})
		return
	}

	score, total := 0, len(quiz.Questions)
	for _, q := range quiz.Questions {
		if ans, ok := input.Answers[fmt.Sprintf("%d", q.ID)]; ok {
			if strings.EqualFold(ans, q.Correct) {
				score++
			}
		}
	}

	percent := 0.0
	if total > 0 {
		percent = float64(score) / float64(total) * 100
	}

	result := models.QuizResult{
		UserID:  c.GetUint("user_id"),
		QuizID:  quiz.ID,
		Score:   score,
		Total:   total,
		Percent: percent,
	}
	if err := database.DB.Create(&result).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Quiz sonucu kaydedilemedi"})
		return
	}

	message := "Maalesef, tekrar deneyin."
	if percent >= 70 {
		message = "Tebrikler, geçtiniz!"
	}
	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"data": gin.H{
			"score":   score,
			"total":   total,
			"percent": percent,
		},
	})
	return
}
