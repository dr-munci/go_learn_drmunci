package database

import (
	"log"

	"golearn/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error
	DB, err = gorm.Open(sqlite.Open("golearn.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Veritabanı bağlantısı başarısız:", err)
	}
	DB.AutoMigrate(
		&models.User{},
		&models.Course{},
		&models.Lesson{},
		&models.Quiz{},
		&models.Question{},
		&models.Progress{},
		&models.QuizResult{},
	)
	log.Println("✅ Veritabanı bağlandı ve tablolar oluşturuldu")
}
