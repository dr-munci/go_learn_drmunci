package handlers

import (
	"net/http"

	"golearn/database"
	"golearn/middleware"
	"golearn/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.Role == "" {
		input.Role = "student"
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hash),
		Role:     input.Role,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Bu e-posta zaten kayıtlı"})
		return
	}
	token, _ := middleware.GenerateToken(user.ID, user.Role)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Kayıt başarılı",
		"token":   token,
		"user":    gin.H{"id": user.ID, "name": user.Name, "email": user.Email, "role": user.Role},
	})
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "E-posta veya şifre hatalı"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "E-posta veya şifre hatalı"})
		return
	}
	token, _ := middleware.GenerateToken(user.ID, user.Role)
	c.JSON(http.StatusOK, gin.H{
		"message": "Giriş başarılı",
		"token":   token,
		"user":    gin.H{"id": user.ID, "name": user.Name, "email": user.Email, "role": user.Role},
	})
}
