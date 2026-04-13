package handlers

import (
	"golearn/database"
	"golearn/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CompleteLesson(c *gin.Context) {
	userID := c.GetUint("user_id")

	var lesson models.Lesson
	if err := database.DB.First(&lesson, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ders bulunamadı"})
		return
	}

	progress := models.Progress{
		UserID:    userID,
		LessonID:  lesson.ID,
		CourseID:  lesson.CourseID,
		Completed: true,
	}
	if err := database.DB.Where("user_id = ? AND lesson_id = ?", userID, lesson.ID).FirstOrCreate(&progress).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "İlerleme kaydedilemedi"})
		return
	}

	progress.Completed = true
	if err := database.DB.Save(&progress).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "İlerleme güncellenemedi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ders tamamlandı"})
	return
}

func GetProgress(c *gin.Context) {
	userID := c.GetUint("user_id")

	type CourseProgress struct {
		CourseID    uint    `json:"course_id"`
		CourseTitle string  `json:"course_title"`
		Total       int     `json:"total_lessons"`
		Completed   int     `json:"completed_lessons"`
		Percent     float64 `json:"percent"`
	}

	var courses []models.Course
	if err := database.DB.Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kurslar getirilemedi"})
		return
	}

	var result []CourseProgress
	for _, course := range courses {
		var total, completed int64
		if err := database.DB.Model(&models.Lesson{}).Where("course_id = ?", course.ID).Count(&total).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ders sayısı hesaplanamadı"})
			return
		}
		if err := database.DB.Model(&models.Progress{}).Where(
			"user_id = ? AND course_id = ? AND completed = ?", userID, course.ID, true,
		).Count(&completed).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "İlerleme hesaplanamadı"})
			return
		}

		if total > 0 {
			result = append(result, CourseProgress{
				CourseID:    course.ID,
				CourseTitle: course.Title,
				Total:       int(total),
				Completed:   int(completed),
				Percent:     float64(completed) / float64(total) * 100,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
	return
}
