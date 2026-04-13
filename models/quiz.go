package models

import "gorm.io/gorm"

type Quiz struct {
	gorm.Model
	Title     string     `json:"title"`
	LessonID  uint       `json:"lesson_id"`
	Questions []Question `json:"questions,omitempty" gorm:"foreignKey:QuizID"`
}

type Question struct {
	gorm.Model
	Text          string `json:"text"`
	OptionA       string `json:"option_a"`
	OptionB       string `json:"option_b"`
	OptionC       string `json:"option_c"`
	OptionD       string `json:"option_d"`
	CorrectAnswer string `json:"-" gorm:"column:correct_answer;not null"`
	QuizID        uint   `json:"quiz_id"`
}

type QuizResult struct {
	gorm.Model
	UserID  uint    `json:"user_id"`
	QuizID  uint    `json:"quiz_id"`
	Score   int     `json:"score"`
	Total   int     `json:"total"`
	Percent float64 `json:"percent"`
}
