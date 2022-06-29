package dbfs

import (
	"crypto/rand"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io"
	"mime"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	ISFILE                   = int8(2)
	ISFOLDER                 = int8(1)
	ISLINK                   = int8(3)
	FILESTATUSWRITING        = int8(0)
	FILESTATUSGOOD           = int8(1)
	FILESTATUSWARNING        = int8(2)
	FILESTATUSOFFLINE        = int8(3)
	FILESTATUSBAD            = int8(4)
	FILESTATUSREBUILDING     = int8(5)
	FILESTATUSDELETED        = int8(6)
	DefaultDataShardsCount   = 5 //EXAMPLE ONLY. REPLACE LATER.
	DefaultParityShardsCount = 2 //EXAMPLE ONLY. REPLACE LATER.
	VERSIONING_OFF           = int8(1)
	VERSIONING_ON_VERSIONS   = int8(2)
	VERSIONING_ON_DELTAS     = int8(3)
)

var (
	ErrFileNotFound     = errors.New("file or folder not found")
	ErrFileFolderExists = errors.New("file/folder already exists")
	ErrFolderNotEmpty   = errors.New("folder contains files")
	ErrNoPermission     = errors.New("no permission")
	ErrNotFolder        = errors.New("not a folder")
	ErrNotFile          = errors.New("not a file")
	ErrInvalidAction    = errors.New("invalid action")
	ErrInvalidFile      = errors.New("invalid file. please check the parameters and try again")
	ErrVersionNotFound  = errors.New("version not found")
)

type File struct {
	FileId             string `gorm:"primaryKey"`
	FileName           string
	MIMEType           string
	EntryType          int8  `gorm:"not null"`
	ParentFolder       *File `gorm:"foreignKey:ParentFolderFileId"`
	ParentFolderFileId string
	VersionNo          int `gorm:"not null"`
	DataId             string
	DataIdVersion      int
	Size               int       `gorm:"not null"`
	ActualSize         int       `gorm:"not null"`
	CreatedTime        time.Time `gorm:"not null"`
	ModifiedUser       *User     `gorm:"foreignKey:ModifiedUserUserId"`
	ModifiedUserUserId string
	ModifiedTime       time.Time `gorm:"not null; autoUpdateTime"`
	VersioningMode     int8      `gorm:"not null"`
	Checksum           string
	FragCount          int
	ParityCount        int
	EncryptionKey      string
	EncryptionIv       string
	PasswordProtected  bool
	LinkFile           *File `gorm:"foreignKey:LinkFileFileId"`
	LinkFileFileId     *string
	LastChecked        time.Time
	Status             int8   `gorm:"not null"`
	HandledServer      string `gorm:"not null"`
}

type FileMetadataModification struct {
	FileName             string
	MIMEType             string
	VersioningMode       int8
	PasswordModification bool
	PasswordProtected    bool
	PasswordHint         string
	OldPassword          string
	NewPassword          string
}

type FileInterface interface {

	// Browse Functions (Local)
	GetFileFragments(tx *gorm.DB, user *User) ([]Fragment, error)
	GetFileMeta(tx *gorm.DB, user *User) error // retrieves all associations (fragments, permissions, etc)
	GetOldVersion(tx *gorm.DB, user *User, versionNo int) (*FileVersion, error)

	// Create Functions (Local)
	CreateSubFolder(tx *gorm.DB, folderName string, user *User) (*File, error)

	// Update Functions (Local)
	UpdateMetaData(tx *gorm.DB, modificationsRequested FileMetadataModification, user *User) error
	rename(tx *gorm.DB, newName string) error
	PasswordProtect(tx *gorm.DB, password string, hint string, user *User) error
	PasswordUnprotect(tx *gorm.DB, password string, user *User) error
	Move(tx *gorm.DB, newParent *File, user *User) error
	Delete(tx *gorm.DB, user *User) error
	AddPermissionUsers(tx *gorm.DB, permission *PermissionNeeded, requestUser *User, users ...User) error
	AddPermissionGroups(tx *gorm.DB, permission *PermissionNeeded, requestUser *User, groups ...Group) error
	RemovePermission(tx *gorm.DB, permission *Permission, user *User) error
	UpdatePermission(tx *gorm.DB, oldPermission *Permission, newPermission *Permission, user *User) error
	UpdateFile(tx *gorm.DB, newSize int, newActualSize int, newEncryptionKey string, checksum string, handlingServer string, user *User) error
	updateFragment(tx *gorm.DB, fragmentId int, fileFragmentPath string, checksum string, serverId string) error
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

	ls, err := ListFilesByFolderId(tx, folderTree[len(folderTree)-1].FileId, user)

	if err != nil {
		return nil, err
	}

	// Finding the file inside
	// Permission check for read is done in ListFilesByFolderId

	destFileName := paths[len(paths)-1]

	for _, file := range ls {
		if file.FileName == destFileName {
			return &file, nil
		}
	}

	return nil, ErrFileNotFound
}

// GetFileById returns a File object based on the FileId given
func GetFileById(tx *gorm.DB, id string, user *User) (*File, error) {
	file := &File{FileId: id}

	err := tx.First(file).Error

	if err != nil {
		return nil, ErrFileNotFound
	}

	hasPermission, err := user.HasPermission(tx, file, &PermissionNeeded{Read: true})

	if !hasPermission {
		return nil, ErrFileNotFound
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

	ls, err := ListFilesByFolderId(tx, folderTree[len(folderTree)-1].FileId, user)

	if err != nil {
		return nil, err
	}

	return ls, nil
}

// ListFilesByFolderId returns an array of File objects based on the FileId/FolderId given
func ListFilesByFolderId(tx *gorm.DB, id string, user *User) ([]File, error) {

	// An easy way to check for permissions.
	_, err := GetFileById(tx, id, user)

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

// transverseByPath Returns an array of FileIds based on the transversal path of pathStringToArray()
// For Example, passing in ["foo", "bar"] will return ["root FileId", "foo FileId", "bar FileId"]
func transverseByPath(tx *gorm.DB, fileNames []string, user *User) ([]File, error) {

	files := make([]File, len(fileNames)+1)

	parentFolder, err := GetRootFolder(tx)

	if err != nil {
		return nil, err
	}

	files[0] = *parentFolder

	for i, fileName := range fileNames {
		fileExists := false
		parentFolderFiles, err := ListFilesByFolderId(tx, parentFolder.FileId, user)

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

	newFolder, err := CreateFolderByParentId(tx, files[len(files)-1].FileId, paths[len(paths)-1], user)

	if err != nil {
		return nil, err
	}

	return newFolder, nil
}

// EXAMPLECreateFile is an example driver for creating a File.
func EXAMPLECreateFile(tx *gorm.DB, user *User, filename string, parentFolderId string) (*File, error) {

	// This is an example script to show how the process should work.

	// First, the system receives the whole file from the user
	// Then, the system creates the record in the system

	file := File{
		FileName:           filename,
		MIMEType:           "",
		ParentFolderFileId: parentFolderId, // root folder for now
		Size:               512,
		VersioningMode:     0,
		Checksum:           "CHECKSUM",
		ParityCount:        5,
		PasswordProtected:  false,
		HandledServer:      "ThisServer",
	}

	err := tx.Transaction(func(tx *gorm.DB) error {

		err := CreateInitialFile(tx, &file, user)
		if err != nil {
			return err
		}

		err = createPermissions(tx, &file)
		if err != nil {
			// By right, there should be no error possible? If any error happens, it's likely a system error.
			// However, in the case there is an error, we will revert the transaction (thus deleting the file entry)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// If no error, file is created. Now the system need to process the file in the pipeline and send it to each server
	// Pipeline can get the amount of parity bits based on the Parity Count, and get the amount of shards based on the amount of servers

	// Then, the system splits the files accordingly

	// Then, the system send each shard to each server, and once it's sent successfully

	for i := 1; i <= int(file.FragCount); i++ {
		fragId := int(i)
		fragmentPath := uuid.New().String()
		serverId := "Server" + strconv.Itoa(i)
		fragChecksum := "CHECKSUM" + strconv.Itoa(i)

		err = createFragment(tx, file.FileId, file.DataId, file.VersionNo, fragId, serverId, fragmentPath, fragChecksum)
		if err != nil {
			// Not sure how to handle this multiple error situation that is possible.
			// Don't necesarrily want to put it in a transaction because I'm worried it'll be too long?
			// or does that make no sense?
			err2 := file.Delete(tx, user)
			if err2 != nil {
				return nil, err2
			}

			return nil, err
			// If fails, delete File record and return error.
		}
	}

	// Now we can update it to indicate that it has been saved successfully.

	err = finishFile(tx, &file, user, "newEncryptionKey", 412)
	if err != nil {
		// If fails, delete File record and return error.
	}

	err = createFileVersionFromFile(tx, &file, user)
	if err != nil {
		// If fails, delete File record and return error.
	}

	return &file, nil

}

// CreateInitialFile if no error, will continue to call create Fragment for each and update
func CreateInitialFile(tx *gorm.DB, file *File, user *User) error {

	// Verify that the file has the following properties
	/*
		1. FileName
		2. Size
		3. ParentFolderFileId
		4. Checksum

		If file.PasswordProtected, ensure that PasswordHint is there. Driver will hold the password.

		If file.ParityCount is not set (or lower than 0), use defaults
	*/

	if file.FileName == "" || file.Size == 0 || file.ParentFolderFileId == "" || file.Checksum == "" {
		return ErrInvalidFile
	}

	dataShardCount := file.FragCount
	if dataShardCount <= 0 {
		// Use default, probably get from REDIS
		// TODO: Replace and implement
		dataShardCount = DefaultDataShardsCount
	}

	parityShardCount := file.ParityCount
	if parityShardCount <= 0 {
		// Use default, probably get from REDIS
		// TODO: Replace and implement
		parityShardCount = DefaultParityShardsCount
	}

	// Verify that the user has permission to the folder they are writing to

	hasPermission, err := user.HasPermission(tx, &File{FileId: file.ParentFolderFileId}, &PermissionNeeded{Read: true})
	if err != nil || !hasPermission {
		return err
	}

	hasPermission, err = user.HasPermission(tx, &File{FileId: file.ParentFolderFileId}, &PermissionNeeded{Write: true})
	if err != nil || !hasPermission {
		return ErrNoPermission
	}

	// Check that there's no file with the same name in the same folder

	ls, err := ListFilesByFolderId(tx, file.ParentFolderFileId, user)

	for _, f := range ls {
		if f.FileName == file.FileName {
			return ErrFileFolderExists
		}
	}

	// Get MIME Type if not present through extension type
	mimetype := file.MIMEType
	if mimetype == "" {
		extension := filepath.Ext(file.FileName)
		if extension != "" {
			mimetype = mime.TypeByExtension(extension)
		}
	}

	newFile := &File{
		FileId:             uuid.New().String(),
		FileName:           file.FileName,
		MIMEType:           mimetype,
		EntryType:          ISFILE,
		ParentFolderFileId: file.ParentFolderFileId,
		VersionNo:          0,
		DataId:             uuid.New().String(),
		DataIdVersion:      0,
		Size:               file.Size,
		ActualSize:         file.Size,
		CreatedTime:        time.Now(),
		ModifiedUserUserId: user.UserId,
		ModifiedTime:       time.Now(),
		VersioningMode:     file.VersioningMode,
		Checksum:           file.Checksum,
		FragCount:          dataShardCount,
		ParityCount:        parityShardCount,
		// Skip EncyrptionKey, input later
		PasswordProtected: file.PasswordProtected,
		LastChecked:       time.Now(),
		Status:            FILESTATUSWRITING,
		HandledServer:     "", // Should be the server Id of the server that is handling the file
	}

	// Create the file
	err = tx.Create(newFile).Error
	if err != nil {
		return err
	}

	*file = *newFile

	// Update and place into FileVersion
	err = createFileVersionFromFile(tx, file, user)
	if err != nil {
		return err
	}

	return nil

}

// createFragment is called whenever the fragment has been written to disk
func createFragment(tx *gorm.DB, fileId string, dataId string, versionNo int, fragId int, serverId string, fragmentPath string, checksum string) error {

	// Validating inputs
	if fileId == "" || dataId == "" || versionNo < 0 || fragId <= 0 || serverId == "" || fragmentPath == "" || checksum == "" {
		return ErrInvalidAction
	}

	newFragment := Fragment{
		FileVersionFileId:    fileId,
		FileVersionDataId:    dataId,
		FileVersionVersionNo: versionNo,
		FragId:               fragId,
		ServerId:             serverId,
		FileFragmentPath:     fragmentPath,
		Checksum:             checksum,
		LastChecked:          time.Now(),
		Status:               FRAGMENTSTATUSGOOD,
	}
	err := tx.Create(&newFragment).Error
	if err != nil {
		return err
	}

	return err
}

// finishFile is called whenever a file has been all written (all fragments written)
func finishFile(tx *gorm.DB, file *File, user *User, newEncryptionKey string, actualSize int) error {

	// Some checks
	if file.Status != FILESTATUSWRITING {
		return ErrInvalidAction
	}
	if actualSize <= 0 {
		return ErrInvalidAction
	}

	// Update the file to be finished
	file.Status = FILESTATUSGOOD
	file.EncryptionKey = newEncryptionKey
	file.ActualSize = actualSize
	file.ModifiedUserUserId = user.UserId
	file.ModifiedTime = time.Now()
	err := tx.Save(file).Error
	if err != nil {
		return err
	}

	// Updating the FileVersion
	err = finaliseFileVersionFromFile(tx, file)
	if err != nil {
		return err
	}

	return nil
}

// CreateFolderByParentId creates a folder based on the id given and returns the folder (File Object)
func CreateFolderByParentId(tx *gorm.DB, id string, folderName string, user *User) (*File, error) {

	// Check that the id exists

	parentFolder := &File{FileId: id}

	err := tx.First(&parentFolder).Error

	if err != nil {
		return nil, ErrFileNotFound
	}

	return parentFolder.CreateSubFolder(tx, folderName, user)

}

// DeleteFileById deletes a file based on the FileId given.
func DeleteFileById(tx *gorm.DB, id string, user *User) error {
	file, err := GetFileById(tx, id, user)
	if err != nil {
		return err
	}

	return file.Delete(tx, user)

}

// DeleteFolderById deletes a folder based on the FileId given.
// Will not delete if there is contents in the folder
func DeleteFolderById(tx *gorm.DB, id string, user *User) error {

	folder, err := GetFileById(tx, id, user)

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

	files, err := ListFilesByFolderId(tx, folder.FileId, user)
	if err != nil {
		return err
	}

	if len(files) > 0 {
		return ErrFolderNotEmpty
	}

	return tx.Delete(folder).Error

}

// DeleteFolderByIdCascade deletes a folder based on the FileId given.
// Will delete all inner contents.
func DeleteFolderByIdCascade(tx *gorm.DB, id string, user *User) error {

	// Checking if the user has permissions

	err := DeleteFolderById(tx, id, user)

	if errors.Is(ErrFolderNotEmpty, err) {
		// Cascading down the folders

		var files []File
		files, err = ListFilesByFolderId(tx, id, user)

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

// deleteSubFoldersCascade - supporter function for DeleteFolderByIdCascade
// Recursively goes through all folders and deletes them.
func deleteSubFoldersCascade(tx *gorm.DB, file *File, user *User) error {

	// Checking for user permission is currently OFF as the way the DB is designed
	// where subdirectories must contain the same permissions at the parent
	// doesn't require us to do so.

	// However, if the design changes, need to implement here.

	// Check if "file" is a file or a folder

	// If "File" is a folder
	if file.EntryType == 0 {
		files, err := ListFilesByFolderId(tx, file.FileId, user)

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
func (f File) GetFileFragments(tx *gorm.DB, user *User) ([]Fragment, error) {
	// Check if the user has read file permissions to it

	hasPermission, err := user.HasPermission(tx, &f, &PermissionNeeded{Read: true})
	if err != nil {
		return nil, err
	} else if !hasPermission {
		return nil, ErrFileNotFound
	}

	var fragments []Fragment
	err = tx.Where("file_version_data_id = ?", f.DataId).Find(&fragments).Error
	if err != nil {
		return nil, err
	}
	return fragments, nil
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
func (f *File) GetOldVersion(tx *gorm.DB, user *User, versionNo int) (*FileVersion, error) {
	// Check if user has read file permissions to it

	hasPermission, err := user.HasPermission(tx, f, &PermissionNeeded{Read: true})
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, ErrFileNotFound
	}

	if versionNo > int(f.VersionNo) {
		return nil, ErrVersionNotFound
	}

	fileVersion, err := getFileVersionFromFile(tx, f, versionNo)
	if err != nil {
		return nil, err
	}

	return fileVersion, nil

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
		folderName, f.FileId).Count(&rows).Error

	if err != nil {
		return nil, err
	}

	if rows >= 1 {
		return nil, ErrFileFolderExists
	}

	newFolder := &File{
		FileId:             uuid.New().String(),
		FileName:           folderName,
		MIMEType:           "",
		EntryType:          0,
		ParentFolder:       f,
		ParentFolderFileId: f.FileId,
		VersionNo:          0,
		DataId:             "",
		DataIdVersion:      0,
		Size:               0,
		ActualSize:         0,
		CreatedTime:        time.Time{},
		ModifiedUserUserId: user.UserId,
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

		err = createPermissions(tx, newFolder)
		if err != nil {
			return err
		}

		err = createFileVersionFromFile(tx, newFolder, user)
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

// UpdateMetaData used to update file's metadata (name, mime type, etc)
func (f *File) UpdateMetaData(tx *gorm.DB, modificationsRequested FileMetadataModification, user *User) error {

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

	if modificationsRequested.PasswordModification {
		// TODO: Check for password if valid

	} else {
		if modificationsRequested.FileName != "" {
			err := f.rename(tx, modificationsRequested.FileName)
			if err != nil {
				return err
			}
		}
		if modificationsRequested.MIMEType != "" {
			f.MIMEType = modificationsRequested.MIMEType
		}
		if modificationsRequested.VersioningMode >= VERSIONING_OFF && modificationsRequested.VersioningMode <= VERSIONING_ON_DELTAS {
			f.VersioningMode = modificationsRequested.VersioningMode
		}
	}

	// increment version number and save it in FileVersion

	f.VersionNo = f.VersionNo + 1

	// Attempt to save and if it fails, rollback

	err = tx.Transaction(func(tx *gorm.DB) error {
		err = createFileVersionFromFile(tx, f, user)
		if err != nil {
			return err
		}

		return tx.Save(f).Error
	})
	return err
}

// rename is a helper function to rename a file for UpdateMetaData
func (f *File) rename(tx *gorm.DB, newName string) error {

	// Check if there is another file that has the same name in the folder

	var rows int64

	err := tx.Model(&File{}).Where("file_name = ? AND parent_folder_file_id = ?",
		newName, f.ParentFolderFileId).Count(&rows).Error

	if err != nil {
		return err
	}

	if rows >= 1 {
		return ErrFileFolderExists
	} else {
		f.FileName = newName
	}
	return nil
}

// PasswordProtect TODO: implement
// Adds a password to a file
func (f *File) PasswordProtect(tx *gorm.DB, password string, hint string, user *User) error {

	// Check if the user has read permission (if not 404)
	hasPermissions, err := user.HasPermission(tx, f, &PermissionNeeded{Read: true})
	if err != nil {
		return err
	} else if !hasPermissions {
		return ErrFileNotFound
	}

	// Check if the user has write permission (if not 403)
	hasPermissions, err = user.HasPermission(tx, f, &PermissionNeeded{Write: true})
	if err != nil {
		return err
	} else if !hasPermissions {
		return ErrNoPermission
	}

	// Create a salt and hash the password

	salt := make([]byte, 32)
	_, err = io.ReadFull(rand.Reader, salt)
	if err != nil {
		return err
	}

	//key, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, 32)
	if err != nil {
		panic(err.Error())
	}

	// Get original encryption key

	//plaintext := f.EncryptionKey

	panic("TeST")
}

// PasswordUnprotect TODO: implement
// Removes a password to a file
func (f File) PasswordUnprotect(tx *gorm.DB, password string, user *User) error {
	//TODO implement me
	panic("implement me")
}

// Move moves a file to a new folder
func (f *File) Move(tx *gorm.DB, newParent *File, user *User) error {

	// Checking that the user has read permissions on the file and the new parent folder
	hasPermissions, err := user.HasPermission(tx, f, &PermissionNeeded{Read: true})
	if err != nil {
		return err
	} else if !hasPermissions {
		return ErrFileNotFound
	}

	hasPermissions, err = user.HasPermission(tx, newParent, &PermissionNeeded{Read: true})
	if err != nil {
		return err
	} else if !hasPermissions {
		return ErrFileNotFound
	}

	// Check that the user has write permissions on the file and the new parent folder

	hasPermissions, err = user.HasPermission(tx, f, &PermissionNeeded{Write: true})
	if err != nil {
		return err
	} else if !hasPermissions {
		return ErrNoPermission
	}

	hasPermissions, err = user.HasPermission(tx, newParent, &PermissionNeeded{Write: true})
	if err != nil {
		return err
	} else if !hasPermissions {
		return ErrNoPermission
	}

	// Update the parent folder of the file
	f.ParentFolderFileId = newParent.FileId

	err = tx.Transaction(func(tx *gorm.DB) error {
		err2 := tx.Save(f).Error
		if err2 != nil {
			return err2
		}

		err2 = createFileVersionFromFile(tx, f, user)

		return err2

	})

	return err

}

// Delete deletes a file or folder and all of its contents
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

	// Delete Versions, Permissions, Delete File

	err = tx.Transaction(func(tx *gorm.DB) error {
		err2 := deleteFileVersionFromFile(tx, f)
		if err2 != nil {
			return err2
		}
		err2 = deleteFilePermissions(tx, f)
		if err != nil {
			return err
		}
		err2 = tx.Delete(f).Error
		return err2

	})

	return err
}

// AddPermissionUsers adds permissions to a file or folder based on a PermissionNeeded struct given.
// takes in file/folder, permission needed, user requesting, and users to apply to
func (f *File) AddPermissionUsers(tx *gorm.DB, permission *PermissionNeeded, requestUser *User, users ...User) error {

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

// AddPermissionGroups adds permissions to a file or folder based on a PermissionNeeded struct given.
// takes in file/folder, permission needed, user requesting, and users to apply to
func (f *File) AddPermissionGroups(tx *gorm.DB, permission *PermissionNeeded, requestUser *User, groups ...Group) error {

	// Checking if the user has permission to share
	permissionsRequiredToShare := permission

	permissionsRequiredToShare.Share = true

	hasPermission, err := requestUser.HasPermission(tx, f, permissionsRequiredToShare)
	if err != nil {
		return err
	}
	if hasPermission {
		err := upsertGroupsPermission(tx, f, permission, requestUser, groups...)
		if err != nil {
			return err
		} else {
			return nil
		}
	} else {
		return ErrNoPermission
	}

}

// RemovePermission revokes permissions for that user or group.
func (f *File) RemovePermission(tx *gorm.DB, permission *Permission, user *User) error {

	var err error

	// Check if it's user or group
	if permission.UserId != "" {
		// User
		err = revokeUsersPermission(tx, f, []User{User{UserId: permission.UserId}}, user)

	} else {
		// Group
		err = revokeGroupsPermission(tx, f, []Group{Group{GroupId: permission.GroupId}})
	}

	return err
}

// UpdatePermission calls upsertUsersPermission or upsertGroupsPermission to update permissions
func (f *File) UpdatePermission(tx *gorm.DB, oldPermission *Permission, newPermission *Permission, user *User) error {

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

	// Change it to a PermissionRequired to work with upsertUsersPermissions.
	// PermissionsNeeded is a struct that contains a boolean for each permission, allows for in theory easier mass setting
	// of permissions.
	PermissionNeeded := &PermissionNeeded{Read: newPermission.CanRead,
		Write: newPermission.CanWrite, Share: newPermission.CanShare,
		Execute: newPermission.CanExecute,
		Audit:   newPermission.Audit}

	// User or Group

	if oldPermission.UserId != "" {
		newUser, err := GetUser(tx, newPermission.UserId)
		if err != nil {
			return err
		}
		err = upsertUsersPermission(tx, f, PermissionNeeded, user, *newUser)
	} else {
		newGroup, err := GetGroupBasedOnGroupId(tx, newPermission.GroupId)
		if err != nil {
			return err
		}
		err = upsertGroupsPermission(tx, f, PermissionNeeded, user, *newGroup)
	}

	return nil
}

func (f *File) GetPermissions(tx *gorm.DB, user *User) ([]Permission, error) {

	// Check if user has read permission (if not 404)
	hasPermissions, err := user.HasPermission(tx, f, &PermissionNeeded{Read: true})
	if err != nil {
		return nil, err
	} else if !hasPermissions {
		return nil, ErrFileNotFound
	}

	var permissions []Permission

	err = tx.Where("file_id = ?", f.FileId).Find(&permissions).Error

	return permissions, err
}

// UpdateFile
// At this stage, the server should have the
// updated file (already processed) , uncompressed file size, compressed data size, and key
func (f *File) UpdateFile(tx *gorm.DB, newSize int, newActualSize int, newEncryptionKey string, checksum string, handlingServer string, user *User) error {

	// Check if user has read permission (if not 404)

	hasPermissions, err := user.HasPermission(tx, f, &PermissionNeeded{Read: true})
	if err != nil {
		return err
	} else if !hasPermissions {
		return ErrFileNotFound
	}

	hasPermissions, err = user.HasPermission(tx, f, &PermissionNeeded{Write: true})
	if err != nil {
		return err
	} else if !hasPermissions {
		return ErrNoPermission
	}

	// Update the file
	f.VersionNo = f.VersionNo + 1
	f.DataId = uuid.New().String()
	f.DataIdVersion = f.DataIdVersion + 1
	f.Size = newSize
	f.ActualSize = newActualSize
	f.ModifiedUserUserId = user.UserId
	f.ModifiedTime = time.Now()
	f.Checksum = checksum
	f.EncryptionKey = newEncryptionKey
	f.LastChecked = time.Now()
	f.Status = FILESTATUSWRITING
	f.HandledServer = handlingServer

	// Updating the file in the database as well as FileVersion

	err = tx.Transaction(func(tx *gorm.DB) error {
		err2 := tx.Save(f).Error
		if err2 != nil {
			return err2
		}
		// Create a new FileVersion
		err2 = createFileVersionFromFile(tx, f, user)
		return err2
	})
	return err
}

// UpdateFragments. Called on each fragment to be created.
func (f *File) updateFragment(tx *gorm.DB, fragmentId int, fileFragmentPath string, checksum string, serverId string) error {
	frag := Fragment{
		FileVersionFileId:        f.FileId,
		FileVersionDataId:        f.DataId,
		FileVersionVersionNo:     f.VersionNo,
		FileVersionDataIdVersion: f.DataIdVersion,
		FragId:                   fragmentId,
		ServerId:                 serverId,
		FileFragmentPath:         fileFragmentPath,
		Checksum:                 checksum,
		LastChecked:              time.Now(),
		FragCount:                f.FragCount,
		Status:                   FRAGMENTSTATUSGOOD,
	}

	return tx.Save(&frag).Error

}

// finishUpdateFile
// Once all fragments are updated, the file is marked as finished.
// Updating Status to FILESTATUSGOOD, and LastChecked to now.
func (f *File) finishUpdateFile(tx *gorm.DB) error {

	err := tx.Transaction(func(tx *gorm.DB) error {
		f.Status = FILESTATUSGOOD
		f.LastChecked = time.Now()

		err2 := tx.Save(f).Error
		if err2 != nil {
			return err2
		}
		// Updating in FileVersion

		// Get FileVersion
		fileVersion := FileVersion{
			FileId:    f.FileId,
			VersionNo: f.VersionNo,
		}

		tx.First(&fileVersion)

		fileVersion.Status = FILESTATUSGOOD
		fileVersion.LastChecked = f.LastChecked

		return tx.Save(&fileVersion).Error

	})
	if err != nil {
		return err
	}

	if f.VersioningMode == VERSIONING_OFF {
		// Mark previous fragment as to be deleted.
		err = tx.Model(&FileVersion{}).Where("file_id = ? AND versioning_mode = ?", f.FileId, VERSIONING_OFF).Update("status", FILESTATUSDELETED).Error

	} else if f.VersioningMode == VERSIONING_ON_VERSIONS {
		err = nil
	} else if f.VersioningMode == VERSIONING_ON_DELTAS {
		panic("Not implemented")
	}

	return err

}

// EXAMPLEUpdateFile
func EXAMPLEUpdateFile(tx *gorm.DB, file *File, user *User) error {

	err := file.UpdateFile(tx, 2020, 2020, "wowKey", "wowChecksum", "wowNew", user)
	if err != nil {
		return err
	}

	// As each fragment is uploaded, each fragment is added to the database.
	for i := 1; i <= int(file.FragCount); i++ {
		err = file.updateFragment(tx, int(i), "wowFragmentPath", "wowChecksum", "wowServerId")
		if err != nil {
			return err
		}
	}

	// Once all fragments are uploaded, the file is marked as finished.
	return file.finishUpdateFile(tx)

}

// ListContents list the contents of a folder (if folder)
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

	err = tx.Model(&File{}).Where("parent_folder_file_id = ?", f.FileId).Find(&files).Error

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
		ls, err := ListFilesByFolderId(tx, f.FileId, user)
		if err != nil {
			return false, err
		} else {
			return len(ls) == 0, nil
		}
	}
}
