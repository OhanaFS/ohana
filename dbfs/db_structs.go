package dbfs

import (
	"gorm.io/gorm"
)

// ! To be implemented

type Server struct {
	gorm.Model
	Name string `gorm:"not null"`
}
