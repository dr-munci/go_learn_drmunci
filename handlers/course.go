package handlers

import (
	"golearn/database"
	"golearn/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetCourses(c *gin.Context) {
	var courses []models.Course
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	category := c.Query("category")
	sort := c.DefaultQuery("sort", "created_at desc")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit
	query := database.DB.Preload("Teacher")
	if category != "" {
		query = query.Where("category = ?", category)
	}

	var total int64
	query.Model(&models.Course{}).Count(&total)
	query.Order(sort).Offset(offset).Limit(limit).Find(&courses)

	c.JSON(http.StatusOK, gin.H{
		"data":  courses,
		"page":  page,
		"limit": limit,
		"total": total,
	})
	return
}

func GetCourse(c *gin.Context) {
	var course models.Course
	if err := database.DB.Preload("Teacher").Preload("Lessons").First(&course, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kurs bulunamadı"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": course})
	return
}

func CreateCourse(c *gin.Context) {
	var course models.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course.TeacherID = c.GetUint("user_id")
	if err := database.DB.Create(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kurs oluşturulamadı"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Kurs oluşturuldu", "data": course})
	return
}

func UpdateCourse(c *gin.Context) {
	var course models.Course
	if err := database.DB.First(&course, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kurs bulunamadı"})
		return
	}

	if course.TeacherID != c.GetUint("user_id") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bu kursu sadece sahibi düzenleyebilir"})
		return
	}

	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kurs güncellenemedi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kurs güncellendi", "data": course})
	return
}

func DeleteCourse(c *gin.Context) {
	var course models.Course
	if err := database.DB.First(&course, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kurs bulunamadı"})
		return
	}

	if course.TeacherID != c.GetUint("user_id") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bu kursu sadece sahibi silebilir"})
		return
	}

	if err := database.DB.Delete(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kurs silinemedi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kurs silindi"})
	return
}
