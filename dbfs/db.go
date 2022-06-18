package dbfs

import (
	"gorm.io/gorm"
	"time"
)

// InitDB Initiates the DB with gorm.db.AutoMigrate
func InitDB(db *gorm.DB) error {
	err := db.AutoMigrate(&User{}, &Group{}, &File{}, &FileVersion{}, &Fragment{}, &Permission{}, &PermissionHistory{})

	if err != nil {
		return err
	}

	// Set root folder

	rootFolder := File{
		FileID:        "00000000-0000-0000-0000-000000000000",
		FileName:      "root",
		EntryType:     0,
		VersionNo:     0,
		Size:          0,
		ActualSize:    0,
		CreatedTime:   time.Time{},
		ModifiedTime:  time.Time{},
		Status:        1,
		HandledServer: "",
	}

	return db.Save(&rootFolder).Error
}

type PermissionNeeded struct {
	Read    bool
	Write   bool
	Execute bool
	Share   bool
}
