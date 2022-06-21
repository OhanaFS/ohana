package dbfs

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
	"time"
)

const (
	ISFILE   = int8(1)
	ISFOLDER = int8(0)
	ISLINK   = int8(2)
)

var (
	ErrFileNotFound     = errors.New("file or folder not found")
	ErrFileFolderExists = errors.New("folder already exists")
	ErrFolderNotEmpty   = errors.New("folder contains files")
	ErrNoPermission     = errors.New("no permission")
	ErrNotFolder        = errors.New("not a folder")
)

type File struct {
	FileID             string `gorm:"primaryKey"`
	FileName           string
	MIMEType           string
	EntryType          int8 `gorm:"not null"`
	ParentFolder       *File
	ParentFolderFileID string
	VersionNo          uint `gorm:"not null"`
	DataID             string
	DataIDVersion      uint
	Size               uint      `gorm:"not null"`
	ActualSize         uint      `gorm:"not null"`
	CreatedTime        time.Time `gorm:"not null"`
	ModifiedUser       User      `gorm:"foreignKey:UserID"`
	ModifiedUserUserID string
	ModifiedTime       time.Time `gorm:"not null; autoUpdateTime"`
	VersioningMode     int8      `gorm:"not null"`
	Checksum           string
	DataShardsCount    uint8
	EncryptionKey      string
	PasswordProtected  bool
	PasswordHint       string
	LinkFileID         *File `gorm:"foreignKey:FileID"`
	LastChecked        time.Time
	Status             int8   `gorm:"not null"`
	HandledServer      string `gorm:"not null"`
}

type FileInterface interface {

	// Browse Functions (Local)
	GetFileFragments(tx *gorm.DB, user *User) ([]Fragment, error)
	GetFileMeta(tx *gorm.DB, user *User) error // retrieves all associations (fragments, permissions, etc)
	GetOldVersion(tx *gorm.DB, user *User, versionNo int) (*FileVersion, error)

	// Create Functions (Local)
	CreateSubFolder(tx *gorm.DB, folderName string, user *User) (*File, error)

	// Update Functions (Local)
	UpdateFragments(tx *gorm.DB, fragments []Fragment, user *User) error
	UpdateMetaData(tx *gorm.DB, file *File, user *User) error
	Rename(tx *gorm.DB, newName string, user *User) error
	PasswordProtect(tx *gorm.DB, password string, hint string, user *User) error
	PasswordUnprotect(tx *gorm.DB, password string, user *User) error
	Move(tx *gorm.DB, newParent *File, user *User) error
	Delete(tx *gorm.DB, user *User) error
	AddPermission(tx *gorm.DB, permission *PermissionNeeded, requestUser *User, users ...User) error
	RemovePermission(tx *gorm.DB, permission *Permission, user *User) error
	UpdatePermission(tx *gorm.DB, oldPermission *Permission, newPermission *Permission, user *User) error
	UpdateFile(tx *gorm.DB, user *User) error
}

var _ FileInterface = &File{}

// GetRootFolder returns the root folder as a File object
func GetRootFolder(tx *gorm.DB) (*File, error) {
	var file File
	err := tx.First(&file, "file_id = ?", "00000000-0000-0000-0000-000000000000").Error

	if err != nil {
		return nil, err
	}

	return &file, nil
}

// GetFileByPath returns a File object based on the path given
func GetFileByPath(tx *gorm.DB, path string, user *User) (*File, error) {

	paths := pathStringToArray(path, true)

	folderTree, err := transverseByPath(tx, paths[0:len(paths)-1], user)

	if err != nil {
		return nil, err
	}

	ls, err := ListFilesByFolderID(tx, folderTree[len(folderTree)-1].FileID, user)

	if err != nil {
		return nil, err
	}

	// Finding the file inside
	// Permission check for read is done in ListFilesByFolderID

	destFileName := paths[len(paths)-1]

	for _, file := range ls {
		if file.FileName == destFileName {
			return &file, nil
		}
	}

	return nil, ErrFileNotFound
}

// GetFileByID returns a File object based on the FileID given
func GetFileByID(tx *gorm.DB, id string, user *User) (*File, error) {
	file := &File{FileID: id}

	err := tx.First(file).Error

	if err != nil {
		return nil, ErrFileNotFound
	}

	hasPermission, err := user.HasPermission(tx, file, &PermissionNeeded{Read: true})

	if !hasPermission {
		return nil, ErrNoPermission
	} else if err != nil {
		return nil, err
	}

	return file, nil

}

// ListFilesByPath returns an array of File objects based on the path given
func ListFilesByPath(tx *gorm.DB, path string, user *User) ([]File, error) {

	paths := pathStringToArray(path, true)

	folderTree, err := transverseByPath(tx, paths, user)

	if err != nil {
		return nil, err
	}

	ls, err := ListFilesByFolderID(tx, folderTree[len(folderTree)-1].FileID, user)

	if err != nil {
		return nil, err
	}

	return ls, nil
}

// ListFilesByFolderID returns an array of File objects based on the FileID/FolderID given
func ListFilesByFolderID(tx *gorm.DB, id string, user *User) ([]File, error) {

	// An easy way to check for permissions.
	_, err := GetFileByID(tx, id, user)

	if err != nil {
		return nil, err
	}

	var files []File

	err = tx.Where("parent_folder_file_id = ?", id).Find(&files).Error

	if err != nil {
		return nil, err
	}

	return files, nil

}

// pathStringToArray returns a string array based on the path string given
// For example, passing in "foo/bar" will return ["foo", "bar"]
func pathStringToArray(path string, fromRoot bool) []string {

	paths := strings.Split(path, "/")

	startNo := 0
	endNo := len(paths) - 1

	if paths[0] == "" && fromRoot { // Does the path start with a '/' ?
		startNo = 1
	}
	if paths[endNo] == "" { // Does the path end with a '/' ?
		endNo = endNo - 1
	}

	return paths[startNo : endNo+1]

}

// transverseByPath Returns an array of FileIDs based on the transversal path of pathStringToArray()
// For Example, passing in ["foo", "bar"] will return ["root FileID", "foo FileID", "bar FileID"]
func transverseByPath(tx *gorm.DB, fileNames []string, user *User) ([]File, error) {

	files := make([]File, len(fileNames)+1)

	parentFolder, err := GetRootFolder(tx)

	if err != nil {
		return nil, err
	}

	files[0] = *parentFolder

	for i, fileName := range fileNames {
		fileExists := false
		parentFolderFiles, err := ListFilesByFolderID(tx, parentFolder.FileID, user)

		if err != nil {
			return nil, err
		}

		for _, file := range parentFolderFiles {
			if file.FileName == fileName {
				fileExists = true
				parentFolder = &file
				break
			}
		}

		if fileExists {
			files[i+1] = *parentFolder
		} else {
			return nil, ErrFileNotFound
		}
	}

	return files, nil

}

// CreateFolderByPath creates a folder based on the path given and returns the folder (File Object)
func CreateFolderByPath(tx *gorm.DB, path string, user *User) (*File, error) {

	paths := pathStringToArray(path, true)

	files, err := transverseByPath(tx, paths[0:(len(paths)-1)], user)

	if err != nil {
		return nil, err
	}

	newFolder, err := CreateFolderByParentID(tx, files[len(files)-1].FileID, paths[len(paths)-1], user)

	if err != nil {
		return nil, err
	}

	return newFolder, nil
}

// CreateFile if no error, will continue to call create Fragment for each and update
func CreateFile(user *User) (*File, error) {
	panic("implement")
}

// CreateFolderByParentID creates a folder based on the id given and returns the folder (File Object)
func CreateFolderByParentID(tx *gorm.DB, id string, folderName string, user *User) (*File, error) {

	// Check that the id exists

	parentFolder := &File{FileID: id}

	err := tx.First(&parentFolder).Error

	if err != nil {
		return nil, ErrFileNotFound
	}

	return parentFolder.CreateSubFolder(tx, folderName, user)

}

// DeleteFileByID deletes a file based on the FileID given.
// TODO: IMPLEMENT
func DeleteFileByID(id string, user *User) error {
	panic("implement")
}

// DeleteFolderByID deletes a folder based on the FileID given.
// Will not delete if there is contents in the folder
func DeleteFolderByID(tx *gorm.DB, id string, user *User) error {

	folder, err := GetFileByID(tx, id, user)

	if err != nil {
		return err
	}

	// Check if user has permissions to delete folder

	hasPermission, err := user.HasPermission(tx, folder, &PermissionNeeded{Write: true})

	if !hasPermission {
		return ErrNoPermission
	} else if err != nil {
		return err
	}

	// Check if the folder has files inside

	files, err := ListFilesByFolderID(tx, folder.FileID, user)
	if err != nil {
		return err
	}

	if len(files) > 0 {
		return ErrFolderNotEmpty
	}

	return tx.Delete(folder).Error

}

// DeleteFolderByIDCascade deletes a folder based on the FileID given.
// Will delete all inner contents.
func DeleteFolderByIDCascade(tx *gorm.DB, id string, user *User) error {

	// Checking if the user has permissions

	err := DeleteFolderByID(tx, id, user)

	if errors.Is(ErrFolderNotEmpty, err) {
		// Cascading down the folders

		var files []File
		files, err = ListFilesByFolderID(tx, id, user)

		if err != nil {
			return err
		}

		err = tx.Transaction(func(tx *gorm.DB) error {

			for _, file := range files {
				err := tx.Transaction(func(tx2 *gorm.DB) error {
					return deleteSubFoldersCascade(tx, &file, user)
				})
				if err != nil {
					return err
				}
			}
			return nil
		})

	} else {
		return err
	}

	return err

}

// deleteSubFoldersCascade - supporter function for DeleteFolderByIDCascade
// Recursively goes through all folders and deletes them.
func deleteSubFoldersCascade(tx *gorm.DB, file *File, user *User) error {

	// Checking for user permission is currently OFF as the way the DB is designed
	// where subdirectories must contain the same permissions at the parent
	// doesn't require us to do so.

	// However, if the design changes, need to implement here.

	// Check if "file" is a file or a folder

	// If "File" is a folder
	if file.EntryType == 0 {
		files, err := ListFilesByFolderID(tx, file.FileID, user)

		if err != nil {
			return err
		}

		for _, file := range files {
			err = tx.Transaction(func(tx *gorm.DB) error {
				return deleteSubFoldersCascade(tx, &file, user)
			})
			if err != nil {
				return err
			}
		}
	}

	if err := deleteFilePermissions(tx, file); err != nil {
		return err
	}

	return tx.Delete(&file).Error

}

// GetFileFragments returns the fragments associated with the File
// TODO: IMPLEMENT
func (f File) GetFileFragments(tx *gorm.DB, user *User) ([]Fragment, error) {
	//TODO implement me
	panic("implement me")
}

// GetFileMeta returns the file metadata including the permissions associated with the File/Folder
func (f *File) GetFileMeta(tx *gorm.DB, user *User) error {

	_, err := user.HasPermission(tx, f, &PermissionNeeded{Read: true})
	if err != nil {
		return ErrFileNotFound
	}

	return tx.Preload(clause.Associations).Find(f).Error
}

// GetOldVersion returns the FileVersion object of the File requested.
// TODO: IMPLEMENT
func (f File) GetOldVersion(tx *gorm.DB, user *User, versionNo int) (*FileVersion, error) {
	//TODO implement me
	panic("implement me")
}

// CreateSubFolder creates a subfolder of the parent folder
func (f *File) CreateSubFolder(tx *gorm.DB, folderName string, user *User) (*File, error) {

	// Check if user has read permission (if not 404)

	hasPermissions, err := user.HasPermission(tx, f, &PermissionNeeded{Read: true})
	if err != nil {
		return nil, err
	} else if !hasPermissions {
		return nil, ErrFileNotFound
	}

	// Check if user has write permission (if not 403)
	hasPermissions, err = user.HasPermission(tx, f, &PermissionNeeded{Write: true})
	if err != nil {
		return nil, err
	} else if !hasPermissions {
		return nil, ErrNoPermission
	}

	var rows int64

	err = tx.Model(&File{}).Where("file_name = ? AND parent_folder_file_id = ?",
		folderName, f.FileID).Count(&rows).Error

	if err != nil {
		return nil, err
	}

	if rows >= 1 {
		return nil, ErrFileFolderExists
	}

	newFolder := &File{
		FileID:             uuid.New().String(),
		FileName:           folderName,
		MIMEType:           "",
		EntryType:          0,
		ParentFolder:       f,
		ParentFolderFileID: f.FileID,
		VersionNo:          0,
		DataID:             "",
		DataIDVersion:      0,
		Size:               0,
		ActualSize:         0,
		CreatedTime:        time.Time{},
		ModifiedUser:       User{},
		ModifiedTime:       time.Time{},
		VersioningMode:     0,
		Status:             1,
		HandledServer:      "",
	}

	// Transaction

	err = tx.Transaction(func(tx *gorm.DB) error {
		err = tx.Save(newFolder).Error

		if err != nil {
			return err
		}

		// TODO: Create Versions

		err = createPermissions(tx, newFolder)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return newFolder, nil
}

func (f File) UpdateFragments(tx *gorm.DB, fragments []Fragment, user *User) error {
	//TODO implement me
	panic("implement me")
}

func (f File) UpdateMetaData(tx *gorm.DB, file *File, user *User) error {
	//TODO implement me
	panic("implement me")
}

// Rename() renames and saves the files instantly .
func (f *File) Rename(tx *gorm.DB, newName string, user *User) error {

	// Check if user has read permission (if not 404)
	hasPermissions, err := user.HasPermission(tx, f, &PermissionNeeded{Read: true})
	if err != nil {
		return err
	} else if !hasPermissions {
		return ErrFileNotFound
	}

	// Check if user has write permission (if not 403)
	hasPermissions, err = user.HasPermission(tx, f, &PermissionNeeded{Write: true})
	if err != nil {
		return err
	} else if !hasPermissions {
		return ErrNoPermission
	}

	// Check if there is another file that has the same name in the folder

	var rows int64

	err = tx.Model(&File{}).Where("file_name = ? AND parent_folder_file_id = ?",
		newName, f.ParentFolderFileID).Count(&rows).Error

	if err != nil {
		return err
	}

	if rows >= 1 {
		return ErrFileFolderExists
	} else {
		f.FileName = newName
		err = tx.Save(f).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (f File) PasswordProtect(tx *gorm.DB, password string, hint string, user *User) error {
	//TODO implement me
	panic("implement me")
}

func (f File) PasswordUnprotect(tx *gorm.DB, password string, user *User) error {
	//TODO implement me
	panic("implement me")
}

func (f File) Move(tx *gorm.DB, newParent *File, user *User) error {
	//TODO implement me
	panic("implement me")
}

// Deletes file or folder and all of its contents

func (f *File) Delete(tx *gorm.DB, user *User) error {

	// Check if user has read permission (if not 404)
	hasPermissions, err := user.HasPermission(tx, f, &PermissionNeeded{Read: true})
	if !hasPermissions {
		return ErrFileNotFound
	} else if err != nil {
		return err
	}

	// Check if user has write permission (if not 403)
	hasPermissions, err = user.HasPermission(tx, f, &PermissionNeeded{Write: true})
	if !hasPermissions {
		return ErrNoPermission
	} else if err != nil {
	}

	// Check if the file is a file or empty folder

	isFileOrEmptyFolder, err := f.IsFileOrEmptyFolder(tx, user)

	if err != nil {
		return err
	}

	if !isFileOrEmptyFolder {
		// Cascading down the folders

		var files []File
		files, err = f.ListContents(tx, user)

		if err != nil {
			return err
		}

		err = tx.Transaction(func(tx *gorm.DB) error {

			for _, file := range files {
				err := tx.Transaction(func(tx2 *gorm.DB) error {
					return deleteSubFoldersCascade(tx, &file, user)
				})
				if err != nil {
					return err
				}
			}
			return nil
		})
	}

	// TODO: Need to handle deletion for versions, fragments etc.
	// Deleting Permissions for the file or folder
	err = deleteFilePermissions(tx, f)
	if err != nil {
		return err
	}

	return tx.Delete(f).Error

}

// AddPermission adds permissions to a file or folder based on a PermissionNeeded struct given.
// takes in file/folder, permission needed, user requesting, and users to apply to
func (f *File) AddPermission(tx *gorm.DB, permission *PermissionNeeded, requestUser *User, users ...User) error {

	// Checking if the user has permission to share
	permissionsRequiredToShare := permission

	permissionsRequiredToShare.Share = true

	hasPermission, err := requestUser.HasPermission(tx, f, permissionsRequiredToShare)
	if err != nil {
		return err
	}
	if hasPermission {
		err := upsertUsersPermission(tx, f, permission, requestUser, users...)
		if err != nil {
			return err
		} else {
			return nil
		}
	} else {
		return ErrNoPermission
	}

}

func (f File) RemovePermission(tx *gorm.DB, permission *Permission, user *User) error {
	//TODO implement me
	panic("implement me")
}

func (f File) UpdatePermission(tx *gorm.DB, oldPermission *Permission, newPermission *Permission, user *User) error {
	//TODO implement me
	panic("implement me")
}

func (f File) UpdateFile(tx *gorm.DB, user *User) error {
	//TODO implement me
	panic("implement me")
}

func (f *File) ListContents(tx *gorm.DB, user *User) ([]File, error) {

	// Check if user has read permission (if not 404)
	hasPermissions, err := user.HasPermission(tx, f, &PermissionNeeded{Read: true})

	if !hasPermissions {
		return nil, ErrFileNotFound
	} else if err != nil {
		return nil, err
	}

	// Check if the file is a folder
	if f.EntryType != ISFOLDER {
		return nil, ErrNotFolder
	}

	var files []File

	err = tx.Where("parent_folder_file_id = ?", f.FileID).Find(&files).Error

	if err != nil {
		return nil, err
	}

	return files, nil

}

// IsFileOrEmptyFolder returns true if the file is a file or an empty folder (useful for permissions)
func (f File) IsFileOrEmptyFolder(tx *gorm.DB, user *User) (bool, error) {
	if f.EntryType == 0 {
		return true, nil
	} else {
		// check that no contents exist
		ls, err := ListFilesByFolderID(tx, f.FileID, user)
		if err != nil {
			return false, err
		} else {
			return len(ls) == 0, nil
		}
	}
}
