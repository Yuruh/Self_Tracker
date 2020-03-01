package database

import (
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
		instance.LogMode(true)
		instance.Set("gorm:auto_preload", true)

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
	instance.AutoMigrate(&models.Connector{})

	//	instance.Create(&models.User{Email: "toto@address.com"})
//	db.Create(&models.User{Email: "tzata@tata.com", Password:"azer"})
}