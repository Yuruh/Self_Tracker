package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
)

type Product struct {
	gorm.Model
	Code string
	Price uint
}

func DataBaseConnect() {
	db, err := gorm.Open("postgres", "user=postgres host=postgres sslmode=disable")
	if err != nil {
		log.Fatalln("failed to connect database", err)
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&Product{})

	// Create
	db.Create(&Product{Code: "L1212", Price: 1000})

	// Read
	var product Product
	db.First(&product, 1) // find product with id 1
	db.First(&product, "code = ?", "L1212") // find product with code l1212

	println(product.Code)

	// Update - update product's price to 2000
	db.Model(&product).Update("Price", 2000)

	// Delete - delete product
	db.Delete(&product)
}