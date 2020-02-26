package models

import (
	"github.com/jinzhu/gorm"
)

// models should be validated outside gorm, using https://github.com/go-playground/validator for example

type User struct {
	gorm.Model
	Email       string  `gorm:"type:varchar(100);unique_index" json:"email"`
	Password	string  `gorm:"not null" json:"password,omitempty"`
	ApiAccess	ApiAccess `gorm:"foreignKey:UserID" json:",omitempty"`
}

type ApiAccess struct {
	gorm.Model
	UserID uint

	// Spotify refresh token
	Spotify string

	// Affect-tag API Key
	AffectTag string

	// Github API ?

	// Steam API ?



	// Do not have public API as of 26/02/2020: Netflix, Origin
}