package dbfs

import (
	"time"
)

type File struct {
	FileID            string `gorm:"primaryKey"`
	FileName          string `gorm:"not null"`
	MIMEType          string
	EntryType         int8  `gorm:"not null"`
	ParentFolder      *File `gorm:"not null"`
	VersionNo         uint  `gorm:"not null"`
	DataID            string
	DataIDVersion     uint
	Size              uint      `gorm:"not null"`
	ActualSize        uint      `gorm:"not null"`
	CreatedTime       time.Time `gorm:"not null"`
	ModifiedUser      User      `gorm:"foreignKey:UserID"`
	ModifiedTime      time.Time `gorm:"not null; autoUpdateTime"`
	VersioningMode    int8      `gorm:"not null"`
	Checksum          string
	DataShardsCount   uint8
	EncryptionKey     string
	PasswordProtected bool
	PasswordHint      string
	LinkFileID        *File `gorm:"foreignKey:FileID"`
	LastChecked       time.Time
	Status            int8   `gorm:"not null"`
	HandledServer     string `gorm:"not null"`
}

type FileInterface interface {

	// Browse Functions (Global)
	GetFileByPath(path string, user User) (*File, error)
	GetFileByID(id string, user User) (*File, error)
	ListFilesByPath(path string, user User) ([]*File, error)
	ListFilesByFolderID(id string, user User) ([]*File, error)

	// Browse Functions (Local)
	GetFileFragments(user User) ([]*Fragment, error)
	GetFileMeta(user User) error // retrieves all associations (fragments, permissions, etc)
	GetOldVersion(user User, versionNo int) (*File, error)

	// Create Functions (Global)
	CreateFolderByPath(path string, user User) (*File, error)
	CreateFile(file *File, user User) (*File, error) // if no error, will continue to call create Fragment for each and update
	CreateFolderByParrentID(id string, folderName string, user User) (*File, error)

	// Create Functions (Local)
	CreateSubFolder(folderName string, user User) (*File, error)

	// Update Functions (Local)
	UpdateFragments(fragments []*Fragment, user User) error
	UpdateMetaData(file *File, user User) error
	Rename(newName string, user User) error
	PasswordProtect(password string, hint string, user User) error
	PasswordUnprotect(password string, user User) error
	Move(newParent *File, user User) error
	Delete(user User) error
	AddPermission(permission *Permission, user User) error
	RemovePermission(permission *Permission, user User) error
	UpdatePermission(oldPermission *Permission, newPermission *Permission, user User) error
	UpdateFile(user User) error

	// Delete Files (Global)
	DeleteFileByID(id string, user User) error
	DeleteFolderByID(id string, user User) error
}
