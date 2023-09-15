package database

import "gorm.io/gorm"

type TokenBlacklist struct {
	gorm.Model
	Token string `gorm:"uniqueIndex"`
}
type User struct {
	gorm.Model
	Username     string `gorm:"uniqueIndex;not null"`
	Password     string `gorm:"not null"`
	RefreshToken string `gorm:"null"`
}
