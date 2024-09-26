package database

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email string `gorm:"uniqueIndex;not null"`
	image string
}

func GetModels() []interface{} {
	return []interface{}{
		&User{},
	}
}
