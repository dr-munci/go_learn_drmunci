package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error
	DB, err = gorm.Open(sqlite.Open("golearn.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Veritabanına bağlanılamadı:", err)
	}
	log.Println("Veritabanı bağlantısı başarılı")
}
