package models

import "gorm.io/gorm"

type Progress struct {
	gorm.Model
	UserID    uint `json:"user_id"`
	LessonID  uint `json:"lesson_id"`
	CourseID  uint `json:"course_id"`
	Completed bool `json:"completed" gorm:"default:false"`
}
