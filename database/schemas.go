package database

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserID uint
	email  string
	image  string
}

func GetModels() []interface{} {
	return []interface{}{
		&User{},
	}
}
