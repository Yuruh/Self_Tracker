package models

import (
	"github.com/jinzhu/gorm"
)

// models should be validated outside gorm, using https://github.com/go-playground/validator for example

type User struct {
	gorm.Model
	Email       string  `gorm:"type:varchar(100);unique_index" json:"email"`
	Password	string  `gorm:"not null" json:"-"`
	Recording	bool	`gorm:"default:false" json:"recording"`
	//	ApiAccess	ApiAccess `gorm:"foreignKey:UserID" json:",omitempty"`

	Connectors []Connector `json:"connectors"`
}
//gorm:"foreignKey:Name"
type Connector struct {
	gorm.Model
	Name		string `json:"name"`
	AvatarUrl	string `json:"avatar_url"`
	Enabled		bool `gorm:"default:false" json:"enabled"`
	Registered	bool `gorm:"default:false" json:"registered"`
	Key			string `json:"-"`
	UserID uint
}

type ApiAccess struct {
	gorm.Model

	// Spotify refresh token

	// Affect-tag API Key

	// Github API ?

	// Steam API ?



	// Do not have public API as of 26/02/2020: Netflix, Origin
}