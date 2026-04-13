package handlers

import (
	"golearn/database"
	"golearn/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetLessons(c *gin.Context) {
	var lessons []models.Lesson
	if err := database.DB.Where("course_id = ?", c.Param("id")).Order(`"order" asc`).Find(&lessons).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Dersler getirilemedi"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": lessons})
	return
}

func CreateLesson(c *gin.Context) {
	var lesson models.Lesson
	if err := c.ShouldBindJSON(&lesson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var course models.Course
	if err := database.DB.First(&course, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kurs bulunamadı"})
		return
	}

	if course.TeacherID != c.GetUint("user_id") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bu kursa sadece sahibi ders ekleyebilir"})
		return
	}

	lesson.CourseID = course.ID
	if err := database.DB.Create(&lesson).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ders oluşturulamadı"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Ders oluşturuldu", "data": lesson})
	return
}

