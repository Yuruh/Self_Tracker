package database

import (
	"fmt"
	"github.com/Yuruh/Self_Tracker/src/database/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"sync"
	"sync/atomic"
)

var mu sync.Mutex
var initialized uint32 = 0
var instance *gorm.DB

func GetDB() *gorm.DB {
	if atomic.LoadUint32(&initialized) == 1 {
		return instance
	}
	mu.Lock()
	defer mu.Unlock()

	if initialized == 0 {
		instance = Connect()
		atomic.StoreUint32(&initialized, 1)
	}

	return instance
}


func Connect() *gorm.DB {
	db, err := gorm.Open("postgres", "user=postgres host=postgres password=changeme sslmode=disable")
	if err != nil {
		log.Fatalln("failed to connect database", err)
	}
	log.Println("Connected to database")
	return db
}

func RunMigration() {
	instance.AutoMigrate(&models.User{})
	instance.AutoMigrate(&models.ApiAccess{})

	var user models.User
	instance.Where("email = ?", "antoine.lempereur@epitech.eu").First(&user)

	var api models.ApiAccess
	if instance.Model(&user).Related(&api).RecordNotFound() {
		fmt.Println("Could not find matching doc")
		instance.Create(&models.ApiAccess{
			UserID: user.ID,
		})
	} else {
		fmt.Println(api.Spotify)
	}



	//	instance.Create(&models.User{Email: "toto@address.com"})
//	db.Create(&models.User{Email: "tzata@tata.com", Password:"azer"})
}