package main

import (
	"golearn/database"
	_ "golearn/docs"
	"golearn/handlers"
	"golearn/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title GoLearn API
// @version 1.0
// @description Uzaktan Egitim Platformu REST API
// @host localhost:8090
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	database.Connect()

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RateLimitMiddleware())

	// Auth routes (herkese açık)
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	// Korumalı route'lar
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/courses", handlers.GetCourses)
		api.GET("/courses/:id", handlers.GetCourse)
		api.POST("/courses", middleware.TeacherOnly(), handlers.CreateCourse)
		api.PUT("/courses/:id", middleware.TeacherOnly(), handlers.UpdateCourse)
		api.DELETE("/courses/:id", middleware.TeacherOnly(), handlers.DeleteCourse)

		api.GET("/courses/:id/lessons", handlers.GetLessons)
		api.POST("/courses/:id/lessons", middleware.TeacherOnly(), handlers.CreateLesson)

		api.GET("/lessons/:id/quiz", handlers.GetQuiz)
		api.POST("/lessons/:id/quiz", middleware.TeacherOnly(), handlers.CreateQuiz)
		api.POST("/quiz/:id/submit", handlers.SubmitQuiz)

		api.GET("/my/progress", handlers.GetProgress)
		api.POST("/lessons/:id/complete", handlers.CompleteLesson)
	}

	r.GET("/ws/classroom/:courseId", middleware.AuthMiddleware(), handlers.ClassroomWS)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8090")
}
