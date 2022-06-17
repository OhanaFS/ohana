package dbfs

import (
	"time"

	"gorm.io/gorm"
)

// ! To be implemented

type Server struct {
	gorm.Model
	Name string `gorm:"not null"`
}

type FileVersion struct {
	FileID            string `gorm:"primaryKey"`
	VersionNo         uint   `gorm:"primaryKey"`
	FileName          string
	MIMEType          string
	EntryType         int8
	ParentFolder      *FileVersion `gorm:"foreignKey:FileID,VersionNo"`
	DataID            string
	DataIDVersion     uint
	Size              uint
	ActualSize        uint
	CreatedTime       time.Time
	ModifiedUser      User `gorm:"foreignKey:UserID;references:FileID,VersionNo"`
	ModifiedTime      time.Time
	VersioningMode    int8
	Checksum          string
	DataShardsCount   uint8
	EncryptionKey     string
	PasswordProtected bool
	PasswordHint      string
	LinkFileID        *FileVersion `gorm:"foreignKey:FileID,VersionNo"`
	LastChecked       time.Time
	Status            int8
	HandledServer     string
	Patch             bool
	PatchBaseVersion  uint
}

// DataID needs  to reference Fragments

type Fragment struct {
	gorm.Model
	DataID           FileVersion `gorm:"primaryKey;foreignKey:FileID,VersionNo"`
	FragID           uint8       `gorm:"primaryKey"`
	ServerID         string
	FileFragmentPath string
	Checksum         string
	LastChecked      time.Time
	Status           int8
}

// Fragment should reference FileVersion

type Permission struct {
	FileID     string `gorm:"primaryKey"`
	User       User   `gorm:"foreignKey:UserID"`
	Group      Group  `gorm:"foreignKey:GroupID"`
	CanRead    bool
	CanWrite   bool
	CanExecute bool
	CanShare   bool
	VersionNo  uint
	Audit      bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

type PermissionHistory struct {
	FileID     string `gorm:"primaryKey"`
	VersionNo  uint   `gorm:"primaryKey"`
	User       User   `gorm:"foreignKey:UserID;references:FileID,VersionNo"`
	Group      Group  `gorm:"foreignKey:GroupID;references:FileID,VersionNo"`
	CanRead    bool
	CanWrite   bool
	CanExecute bool
	CanShare   bool
	Audit      bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}
