package database

import "gorm.io/gorm"

type User struct {
	gorm.Model
	GoogleID string `gorm:"uniqueIndex;not null"`
	Email    string `gorm:"uniqueIndex;not null"`
	Image    string
}

func GetModels() []interface{} {
	return []interface{}{
		&User{},
	}
}
