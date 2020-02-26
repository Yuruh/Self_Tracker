package models

import (
	"github.com/jinzhu/gorm"
)

// models should be validated outside gorm, using https://github.com/go-playground/validator for example

type User struct {
	gorm.Model
	Email        string  `gorm:"type:varchar(100);unique_index" json:"email"`
	Password 	 string  `gorm:"not null" json:"password,omitempty"`
}
