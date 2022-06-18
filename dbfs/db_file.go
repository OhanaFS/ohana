package dbfs

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

var (
	ErrFileNotFound = errors.New("file or folder not found")
)

type File struct {
	FileID            string `gorm:"primaryKey"`
	FileName          string
	MIMEType          string
	EntryType         int8 `gorm:"not null"`
	ParentFolder      *File
	VersionNo         uint `gorm:"not null"`
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
	GetRootFolder() (*File, error)
	GetFileByPath(path string, user User) (*File, error)
	GetFileByID(id string, user User) (*File, error)
	ListFilesByPath(path string, user User) (*[]File, error)
	ListFilesByFolderID(id string, user User) (*[]File, error)

	// Browse Functions (Local)
	GetFileFragments(user User) ([]*Fragment, error)
	GetFileMeta(user User) error // retrieves all associations (fragments, permissions, etc)
	GetOldVersion(user User, versionNo int) (*File, error)

	// Create Functions (Global)
	CreateFolderByPath(path string, user User) (*File, error)
	CreateFile(user User) (*File, error) // if no error, will continue to call create Fragment for each and update
	CreateFolderByParrentID(id string, folderName string, user User) (*File, error)

	// Create Functions (Local)
	CreateSubFolder(folderName string, user User) (*File, error)

	// Update Functions (Local)
	UpdateFragments(fragments *[]Fragment, user User) error
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

// path is like

func GetRootFolder(tx *gorm.DB) (*File, error) {
	var file File
	err := tx.First(&file, "file_id = ?", "00000000-0000-0000-0000-000000000000").Error

	if err != nil {
		return nil, err
	}

	return &file, nil
}

func ListFilesByFolderID(tx *gorm.DB, id string, user User) (*[]File, error) {

	var files []File

	err := tx.Where("parent_folder = ?", id).Find(&files).Error

	if err != nil {
		return nil, err
	}
	return &files, nil

}

func ListFilesByPath(path string, user User) (*[]File, error) {

	// If GetFileByPath works use similar code.

	return nil, nil
}

func GetFileByPath(tx *gorm.DB, path string, user User) (*File, error) {

	// ! Need to implement permissions check

	paths := strings.Split(path, "/")

	startNo := 0
	endNo := len(paths) - 1

	if paths[0] == "" { // Does the path start with a '/' ?
		startNo = 1
	}
	if paths[endNo] == "" { // Does the path end with a '/' ?
		endNo = endNo - 1
	}

	parentFolder, err := GetRootFolder(tx)
	parentFolderFiles, err := ListFilesByFolderID(tx, parentFolder.FileID, user)

	if err != nil {
		return nil, err
	}

	// Transversing the tree
	for i := startNo; i <= endNo; i++ {

		// Checking if the file exists
		fileExists := false

		for _, file := range *parentFolderFiles {
			if file.FileName == paths[i] {
				fileExists = true
				parentFolder = &file
				break
			}
		}

		if fileExists && i < endNo {
			parentFolderFiles, err = ListFilesByFolderID(tx, parentFolder.FileID, user)
		} else if fileExists && i == endNo {
			break
		} else {
			return nil, ErrFileNotFound
		}

	}

	return parentFolder, nil
}

func GetFileByID(tx *gorm.DB, id string, user User) (*File, error) {
	file := &File{FileID: id}

	err := tx.First(file).Error

	if err != nil {
		return nil, ErrFileNotFound
	}

	err = user.HasPermission(tx, file, PermissionNeeded{Read: true})

	if err != nil {
		// Check if permission invalid, if so return not found, else return error
		return nil, ErrFileNotFound
	}

	return file, nil

}

func CreateFolderByParrentID(tx *gorm.DB, id string, folderName string, user User) (*File, error) {

	// Check tha the id exists

	parentFolder := &File{FileID: id}

	err := tx.First(parentFolder).Error

	if err != nil {
		return nil, ErrFileNotFound
	}

	// Check if the user has the right permissions

	/*
		err = user.HasPermission(tx, parentFolder, PermissionNeeded{Read: true})

		if err != nil{
			if errors.Is(err, ErrNoPermission){
				return ErrFileNotFound
			}
		}

		err = user.HasPermission(tx, parentFolder, PermissionNeeded{Write: true})

	*/

	newFolder := &File{
		FileID:         uuid.New().String(),
		FileName:       folderName,
		MIMEType:       "",
		EntryType:      0,
		ParentFolder:   parentFolder,
		VersionNo:      0,
		DataID:         "",
		DataIDVersion:  0,
		Size:           0,
		ActualSize:     0,
		CreatedTime:    time.Time{},
		ModifiedUser:   User{},
		ModifiedTime:   time.Time{},
		VersioningMode: 0,
		Status:         1,
		HandledServer:  "",
	}

	err = tx.Save(newFolder).Error

	if err != nil {
		return nil, err
	}

	// Create Versions

	// Create Permissions

	return newFolder, nil

}
