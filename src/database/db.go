package database

import (
	"github.com/Yuruh/Self_Tracker/src/database/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
)

func Connect() *gorm.DB {
	db, err := gorm.Open("postgres", "user=postgres host=postgres password=changeme sslmode=disable")
	if err != nil {
		log.Fatalln("failed to connect database", err)
	}
	log.Println("Connected to database")
	return db
}

func RunMigration(db *gorm.DB) {
	db.AutoMigrate(&models.User{})

//	db.Create(&models.User{Email: "toto@address.com"})
//	db.Create(&models.User{Email: "tzata@tata.com", Password:"azer"})

	var user models.User
	db.First(&user, 1)

	println(user.Email)
}