package db

import (
	"gorm.io/gorm"
)

type Shortener struct {
	gorm.Model
	Link    string `gorm:"not null;unique"`
	NewLink string `gorm:"type:varchar(10);not null"`
	Count   uint   `gorm:"not null"`
}
