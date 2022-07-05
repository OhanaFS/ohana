package dbfs

import (
	"gorm.io/gorm"
)

// ! To be implemented

type Server struct {
	gorm.Model
	Name string `gorm:"not null"`
}

type Role struct {
	RoleID      int `gorm:"primaryKey;autoIncrement"`
	RoleMapping string
	RoleName    string   `gorm:"not null"`
	Users       []*User  `gorm:"many2many:user_roles;"`
	Groups      []*Group `gorm:"many2many:group_roles;"`
}
