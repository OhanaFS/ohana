package dbfs

import (
	"gorm.io/gorm"
)

// ! To be implemented

type Server struct {
	gorm.Model
	Name string `gorm:"not null"`
}

type PasswordProtect struct {
	FileId       string `gorm:"primaryKey"`
	PasswordSalt string
	PasswordIv   string
	PasswordHint string
}
