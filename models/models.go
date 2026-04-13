package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`
	FullName string `gorm:"not null" json:"full_name"`
	Role     string `gorm:"default:student" json:"role"` // student, teacher, admin
}

type Course struct {
	gorm.Model
	Title       string   `gorm:"not null" json:"title"`
	Description string   `json:"description"`
	TeacherID   uint     `gorm:"not null" json:"teacher_id"`
	Teacher     User     `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
	Lessons     []Lesson `json:"lessons,omitempty"`
}

type Lesson struct {
	gorm.Model
	Title    string `gorm:"not null" json:"title"`
	Content  string `json:"content"`
	VideoURL string `json:"video_url"`
	OrderNum int    `gorm:"not null" json:"order_num"`
	CourseID uint   `gorm:"not null" json:"course_id"`
	Course   Course `gorm:"foreignKey:CourseID" json:"-"`
}

type Enrollment struct {
	gorm.Model
	UserID   uint   `gorm:"not null" json:"user_id"`
	CourseID uint   `gorm:"not null" json:"course_id"`
	User     User   `gorm:"foreignKey:UserID" json:"-"`
	Course   Course `gorm:"foreignKey:CourseID" json:"-"`
}

type Quiz struct {
	gorm.Model
	LessonID  uint       `gorm:"not null" json:"lesson_id"`
	Lesson    Lesson     `gorm:"foreignKey:LessonID" json:"-"`
	Questions []Question `json:"questions,omitempty"`
}

type Question struct {
	gorm.Model
	QuizID        uint   `gorm:"not null" json:"quiz_id"`
	Text          string `gorm:"not null" json:"text"`
	OptionA       string `gorm:"not null" json:"option_a"`
	OptionB       string `gorm:"not null" json:"option_b"`
	OptionC       string `gorm:"not null" json:"option_c"`
	OptionD       string `gorm:"not null" json:"option_d"`
	CorrectAnswer string `gorm:"not null" json:"-"`
}

type QuizResult struct {
	gorm.Model
	UserID      uint      `gorm:"not null" json:"user_id"`
	QuizID      uint      `gorm:"not null" json:"quiz_id"`
	Score       int       `gorm:"not null" json:"score"`
	TotalPoints int       `gorm:"not null" json:"total_points"`
	CompletedAt time.Time `json:"completed_at"`
	User        User      `gorm:"foreignKey:UserID" json:"-"`
	Quiz        Quiz      `gorm:"foreignKey:QuizID" json:"-"`
}

type LessonProgress struct {
	gorm.Model
	UserID      uint      `gorm:"not null" json:"user_id"`
	LessonID    uint      `gorm:"not null" json:"lesson_id"`
	Completed   bool      `gorm:"default:false" json:"completed"`
	CompletedAt time.Time `json:"completed_at"`
	User        User      `gorm:"foreignKey:UserID" json:"-"`
	Lesson      Lesson    `gorm:"foreignKey:LessonID" json:"-"`
}
