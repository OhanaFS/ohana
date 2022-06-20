package dbfs

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

var (
	ErrFileNotFound   = errors.New("file or folder not found")
	ErrFolderExists   = errors.New("folder already exists")
	ErrFolderNotEmpty = errors.New("folder contains files")
	ErrNoPermission   = errors.New("no permission")
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
	GetFileFragments(tx *gorm.DB, user User) ([]Fragment, error)
	GetFileMeta(tx *gorm.DB, user User) error // retrieves all associations (fragments, permissions, etc)
	GetOldVersion(tx *gorm.DB, user User, versionNo int) (*FileVersion, error)

	// Create Functions (Local)
	CreateSubFolder(tx *gorm.DB, folderName string, user User) (*File, error)

	// Update Functions (Local)
	UpdateFragments(tx *gorm.DB, fragments []Fragment, user User) error
	UpdateMetaData(tx *gorm.DB, file *File, user User) error
	Rename(tx *gorm.DB, newName string, user User) error
	PasswordProtect(tx *gorm.DB, password string, hint string, user User) error
	PasswordUnprotect(tx *gorm.DB, password string, user User) error
	Move(tx *gorm.DB, newParent *File, user User) error
	Delete(tx *gorm.DB, user User) error
	AddPermission(tx *gorm.DB, permission PermissionNeeded, requestUser User, users ...User) error
	RemovePermission(tx *gorm.DB, permission *Permission, user User) error
	UpdatePermission(tx *gorm.DB, oldPermission *Permission, newPermission *Permission, user User) error
	UpdateFile(tx *gorm.DB, user User) error
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
func GetFileByPath(tx *gorm.DB, path string, user User) (*File, error) {

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
func GetFileByID(tx *gorm.DB, id string, user User) (*File, error) {
	file := &File{FileID: id}

	err := tx.First(file).Error

	if err != nil {
		return nil, ErrFileNotFound
	}

	//err = user.HasPermission(tx, file, PermissionNeeded{Read: true})

	if err != nil {
		// Check if permission invalid, if so return not found, else return error
		return nil, ErrFileNotFound
	}

	return file, nil

}

// ListFilesByPath returns an array of File objects based on the path given
func ListFilesByPath(tx *gorm.DB, path string, user User) ([]File, error) {

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
func ListFilesByFolderID(tx *gorm.DB, id string, user User) ([]File, error) {

	var files []File

	err := tx.Where("parent_folder_file_id = ?", id).Find(&files).Error

	if err != nil {
		return nil, err
	}

	// TODO: Check permissions

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
func transverseByPath(tx *gorm.DB, fileNames []string, user User) ([]File, error) {

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

			// TODO: Check if user has permission

			files[i+1] = *parentFolder
		} else {
			return nil, ErrFileNotFound
		}
	}

	return files, nil

}

// CreateFolderByPath creates a folder based on the path given and returns the folder (File Object)
func CreateFolderByPath(tx *gorm.DB, path string, user User) (*File, error) {

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
func CreateFile(user User) (*File, error) {
	panic("implement")
}

// CreateFolderByParentID creates a folder based on the id given and returns the folder (File Object)
// TODO: Move the implementation here into an local function
func CreateFolderByParentID(tx *gorm.DB, id string, folderName string, user User) (*File, error) {

	// Check that the id exists

	parentFolder := &File{FileID: id}

	err := tx.First(&parentFolder).Error

	if err != nil {
		return nil, ErrFileNotFound
	}

	// Check if the user has read permissions

	hasPermissions, err := user.HasPermission(tx, parentFolder, PermissionNeeded{Read: true})

	if err != nil {
		return nil, err
	} else if !hasPermissions {
		return nil, ErrFileNotFound
	}

	// Check if the user has write permissions
	hasPermissions, err = user.HasPermission(tx, parentFolder, PermissionNeeded{Write: true})

	if err != nil {
		return nil, err
	} else if !hasPermissions {
		return nil, ErrNoPermission
	}

	// Check that folder doesn't exist

	var rows int64

	err = tx.Model(&File{}).Where("file_name = ? AND parent_folder_file_id = ?",
		folderName, parentFolder.FileID).Count(&rows).Error

	if err != nil {
		return nil, err
	}

	if rows >= 1 {
		return nil, ErrFolderExists
	}

	newFolder := &File{
		FileID:             uuid.New().String(),
		FileName:           folderName,
		MIMEType:           "",
		EntryType:          0,
		ParentFolder:       parentFolder,
		ParentFolderFileID: parentFolder.FileID,
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

// DeleteFileByID deletes a file based on the FileID given.
// TODO: IMPLEMENT
func DeleteFileByID(id string, user User) error {
	panic("implement")
}

// DeleteFolderByID deletes a folder based on the FileID given.
// Will not delete if there is contents in the folder
func DeleteFolderByID(tx *gorm.DB, id string, user User) error {

	folder, err := GetFileByID(tx, id, user)

	if err != nil {
		return err
	}

	// Check if user has permissions to delete folder

	//err = user.HasPermission(tx, folder, PermissionNeeded{Write: true})
	if err != nil {
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
func DeleteFolderByIDCascade(tx *gorm.DB, id string, user User) error {

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
					return deleteSubFoldersCascade(tx, file, user)
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
func deleteSubFoldersCascade(tx *gorm.DB, file File, user User) error {

	// Checking for user permission is currently OFF as the way the DB is designed
	// where subdirectories must contain the same permissions at the parent
	// doesn't require us to do so.

	// However, if the design changes, need to implement here.

	// Check if "file" is a file or a folder

	var err error

	err = nil

	// If "File" is a folder
	if file.EntryType == 0 {
		files, err := ListFilesByFolderID(tx, file.FileID, user)

		if err != nil {
			return err
		}

		for _, file := range files {
			err = tx.Transaction(func(tx *gorm.DB) error {
				return deleteSubFoldersCascade(tx, file, user)
			})
			if err != nil {
				return err
			}
		}

	}

	err = file.Delete(tx, user)

	return err

}

// GetFileFragments returns the fragments associated with the File
// TODO: IMPLEMENT
func (f File) GetFileFragments(tx *gorm.DB, user User) ([]Fragment, error) {
	//TODO implement me
	panic("implement me")
}

// GetFileMeta returns the file metadata including the permissions associated with the File/Folder
// TODO: IMPLEMENT
func (f File) GetFileMeta(tx *gorm.DB, user User) error {
	//TODO implement me
	panic("implement me")
}

// GetOldVersion returns the FileVersion object of the File requested.
// TODO: IMPLEMENT
func (f File) GetOldVersion(tx *gorm.DB, user User, versionNo int) (*FileVersion, error) {
	//TODO implement me
	panic("implement me")
}

func (f File) CreateSubFolder(tx *gorm.DB, folderName string, user User) (*File, error) {
	//TODO implement me
	panic("implement me")
}

func (f File) UpdateFragments(tx *gorm.DB, fragments []Fragment, user User) error {
	//TODO implement me
	panic("implement me")
}

func (f File) UpdateMetaData(tx *gorm.DB, file *File, user User) error {
	//TODO implement me
	panic("implement me")
}

func (f File) Rename(tx *gorm.DB, newName string, user User) error {
	//TODO implement me
	panic("implement me")
}

func (f File) PasswordProtect(tx *gorm.DB, password string, hint string, user User) error {
	//TODO implement me
	panic("implement me")
}

func (f File) PasswordUnprotect(tx *gorm.DB, password string, user User) error {
	//TODO implement me
	panic("implement me")
}

func (f File) Move(tx *gorm.DB, newParent *File, user User) error {
	//TODO implement me
	panic("implement me")
}

func (f File) Delete(tx *gorm.DB, user User) error {
	//TODO NOT FULLY IMPLEMENTED. NO CHECKS DONE
	return tx.Delete(f).Error

}

// AddPermission adds permissions to a file or folder based on a PermissionNeeded struct given.
// takes in file/folder, permission needed, user requesting, and users to apply to
func (f *File) AddPermission(tx *gorm.DB, permission PermissionNeeded, requestUser User, users ...User) error {

	// Checking if the user has permission to share
	permissionsRequiredToShare := permission

	permissionsRequiredToShare.Share = true

	hasPermission, err := requestUser.HasPermission(tx, f, permissionsRequiredToShare)
	if err != nil {
		return err
	}
	if hasPermission {
		err := upsertUsersPermission(tx, f, permission, users...)
		if err != nil {
			return err
		} else {
			return nil
		}
	} else {
		return ErrNoPermission
	}

}

func (f File) RemovePermission(tx *gorm.DB, permission *Permission, user User) error {
	//TODO implement me
	panic("implement me")
}

func (f File) UpdatePermission(tx *gorm.DB, oldPermission *Permission, newPermission *Permission, user User) error {
	//TODO implement me
	panic("implement me")
}

func (f File) UpdateFile(tx *gorm.DB, user User) error {
	//TODO implement me
	panic("implement me")
}

// IsFileOrEmptyFolder returns true if the file is a file or an empty folder (useful for permissions)
func (f File) IsFileOrEmptyFolder(tx *gorm.DB, user User) (bool, error) {
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
