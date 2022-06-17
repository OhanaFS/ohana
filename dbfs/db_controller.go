package dbfs

import (
	"gorm.io/gorm"
)

// InitDB Initiates the DB with gorm.db.AutoMigrate
func InitDB(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &Group{}, &File{}, &FileVersion{}, &Fragment{}, &Permission{}, &PermissionHistory{})
}

type PermissionNeeded struct {
	Read    bool
	Write   bool
	Execute bool
	Share   bool
}
