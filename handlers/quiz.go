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
	type createQuestionInput struct {
		Text    string `json:"text" binding:"required"`
		OptionA string `json:"option_a" binding:"required"`
		OptionB string `json:"option_b" binding:"required"`
		OptionC string `json:"option_c" binding:"required"`
		OptionD string `json:"option_d" binding:"required"`
		Correct string `json:"correct" binding:"required,oneof=a b c d A B C D"`
	}
	type createQuizInput struct {
		Title     string                `json:"title" binding:"required"`
		Questions []createQuestionInput `json:"questions" binding:"required,min=1"`
	}

	var input createQuizInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var lesson models.Lesson
	if err := database.DB.First(&lesson, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ders bulunamadı"})
		return
	}

	quiz := models.Quiz{
		Title:    input.Title,
		LessonID: lesson.ID,
	}
	for _, q := range input.Questions {
		quiz.Questions = append(quiz.Questions, models.Question{
			Text:          q.Text,
			OptionA:       q.OptionA,
			OptionB:       q.OptionB,
			OptionC:       q.OptionC,
			OptionD:       q.OptionD,
			CorrectAnswer: strings.ToLower(q.Correct),
		})
	}

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
			if strings.EqualFold(ans, q.CorrectAnswer) {
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
