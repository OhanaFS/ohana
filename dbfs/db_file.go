package dbfs

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"golang.org/x/crypto/scrypt"
	"io"
	"mime"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/OhanaFS/ohana/util/random_phrase_generator"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	IsFile                   = int8(2)
	IsFolder                 = int8(1)
	IsLink                   = int8(3)
	FileStatusWriting        = int8(0)
	FileStatusGood           = int8(1)
	FileStatusWarning        = int8(2)
	FileStatusOffline        = int8(3)
	FileStatusBad            = int8(4)
	FileStatusReBuilding     = int8(5)
	FileStatusToBeDeleted    = int8(6)
	FileStatusDeleted        = int8(7)
	DefaultDataShardsCount   = 5 //EXAMPLE ONLY. REPLACE LATER.
	DefaultParityShardsCount = 2 //EXAMPLE ONLY. REPLACE LATER.
	VersioningOff            = int8(1)
	VersioningOnVersions     = int8(2)
	VersioningOnDeltas       = int8(3)
	UsersFolderId            = "00000000-0000-0000-0000-000000000001"
	GroupsFolderId           = "00000000-0000-0000-0000-000000000002"
	SharedURLRegexConstriant = "^[a-zA-Z0-9]{1,40}$"
)

var (
	ErrFileNotFound       = errors.New("file or folder not found")
	ErrFileFolderExists   = errors.New("file/folder already exists")
	ErrFolderNotEmpty     = errors.New("folder contains files")
	ErrNoPermission       = errors.New("no permission")
	ErrNotFolder          = errors.New("not a folder")
	ErrNotFile            = errors.New("not a file")
	ErrInvalidAction      = errors.New("invalid action")
	ErrInvalidFile        = errors.New("invalid file. please check the parameters and try again")
	ErrVersionNotFound    = errors.New("version not found")
	ErrPasswordRequired   = errors.New("password required")
	ErrNoPassword         = errors.New("no password")
	ErrLinkExists         = errors.New("link already exists")
	ErrSharedLinkNotFound = errors.New("shared link not found")
	ErrLinkConstraint     = errors.New("link must be alphanumerical and between 1-40 characters")
)

type File struct {
	FileId             string    `gorm:"primaryKey" json:"file_id"`
	FileName           string    `json:"file_name"`
	MIMEType           string    `json:"mime_type"`
	EntryType          int8      `gorm:"not null" json:"entry_type"`
	ParentFolder       *File     `gorm:"foreignKey:ParentFolderFileId" json:"-"`
	ParentFolderFileId *string   `json:"parent_folder_id"`
	VersionNo          int       `gorm:"not null" json:"version_no"`
	DataId             string    `json:"-"` //TODO: Convert this to a pointer to a string and make it unique or nullable
	DataIdVersion      int       `json:"data_version_no"`
	Size               int64     `gorm:"not null" json:"size"`
	ActualSize         int64     `gorm:"not null" json:"actual_size"`
	CreatedTime        time.Time `gorm:"not null"  json:"created_time"`
	ModifiedUser       *User     `gorm:"foreignKey:ModifiedUserUserId" json:"-"`
	ModifiedUserUserId *string   `json:"modified_user_user_id"`
	ModifiedTime       time.Time `gorm:"not null; autoUpdateTime" json:"modified_time"`
	VersioningMode     int8      `gorm:"not null" json:"versioning_mode"`
	Checksum           string    `json:"checksum"`
	TotalShards        int       `json:"total_shards"`
	DataShards         int       `json:"data_shards"`
	ParityShards       int       `json:"parity_shards"`
	KeyThreshold       int       `json:"key_threshold"`
	EncryptionKey      string    `json:"-"`
	EncryptionIv       string    `json:"-"`
	PasswordProtected  bool      `json:"password_protected"`
	LinkFile           *File     `gorm:"foreignKey:LinkFileFileId" json:"-"`
	LinkFileFileId     *string   `json:"link_file_id"`
	LastChecked        time.Time `json:"last_checked"`
	Status             int8      `gorm:"not null" json:"status"`
	HandledServer      string    `gorm:"not null" json:"-"`
}

type FileMetadataModification struct {
	FileName             string `json:"file_name"`
	MIMEType             string `json:"mime_type"`
	VersioningMode       int8   `json:"versioning_mode"`
	PasswordModification bool   `json:"password_modification"`
	PasswordProtected    bool   `json:"password_protected"`
	PasswordHint         string `json:"password_hint"`
	OldPassword          string `json:"old_password"`
	NewPassword          string `json:"new_password"`
}

type FileInterface interface {

	// Browse Functions (Local)
	GetFileFragments(tx *gorm.DB, user *User) ([]Fragment, error)
	GetFileMeta(tx *gorm.DB, user *User) error // retrieves all associations (fragments, permissions, etc)
	GetOldVersion(tx *gorm.DB, user *User, versionNo int) (*FileVersion, error)

	// Create Functions (Local)
	CreateSubFolder(tx *gorm.DB, folderName string, user *User, server string) (*File, error)

	// Update Functions (Local)
	UpdateMetaData(tx *gorm.DB, modificationsRequested FileMetadataModification, user *User) error
	rename(tx *gorm.DB, newName string) error
	PasswordProtect(tx *gorm.DB, oldPassword string, newPassword string, hint string, user *User) error
	PasswordUnprotect(tx *gorm.DB, password string, user *User) error
	Move(tx *gorm.DB, newParent *File, user *User) error
	Delete(tx *gorm.DB, user *User, server string) error
	AddPermissionUsers(tx *gorm.DB, permission *PermissionNeeded, requestUser *User, users ...User) error
	AddPermissionGroups(tx *gorm.DB, permission *PermissionNeeded, requestUser *User, groups ...Group) error
	RemovePermission(tx *gorm.DB, permission *Permission, user *User) error
	UpdatePermission(tx *gorm.DB, oldPermission *Permission, newPermission *Permission, user *User) error
	UpdateFile(tx *gorm.DB, newSize int64, newActualSize int64, checksum string, handlingServer string, dataKey string, dataIv string, password string, user *User, newFileName string) error
	UpdateFragment(tx *gorm.DB, fragmentId int, fileFragmentPath string, checksum string, serverId string) error
	CreateSharedLink(tx *gorm.DB, user *User, link string) (*SharedLink, error)
	GetSharedLinks(tx *gorm.DB, user *User) ([]SharedLink, error)
	DeleteSharedLink(tx *gorm.DB, user *User, link string) error
	UpdateSharedLink(tx *gorm.DB, user *User, link string, newLink string) error
	AddToFavorites(tx *gorm.DB, user *User) error
	RemoveFromFavorites(tx *gorm.DB, user *User) error
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

// GetHomeFolder Returns the user's home folder as a File Object.
func GetHomeFolder(tx *gorm.DB, user *User) (*File, error) {

	file, err := GetFileById(tx, user.HomeFolderId, user)
	if err != nil {
		return nil, err
	}

	return file, nil

}

// LsHomeFolder returns the contents of the user's home folder as a list of File objects
func LsHomeFolder(tx *gorm.DB, user *User) ([]File, error) {
	return ListFilesByFolderId(tx, user.HomeFolderId, user)
}

// GetFileByPath returns a File object based on the path given
func GetFileByPath(tx *gorm.DB, path string, user *User, fromHome bool) (*File, error) {

	paths := pathStringToArray(path, true)

	folderTree, err := traverseByPath(tx, paths[0:len(paths)-1], user, fromHome)

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
func ListFilesByPath(tx *gorm.DB, path string, user *User, fromHome bool) ([]File, error) {

	paths := pathStringToArray(path, true)

	folderTree, err := traverseByPath(tx, paths, user, fromHome)

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

// traverseByPath Returns an array of FileIds based on the transversal path of pathStringToArray()
// For Example, passing in ["foo", "bar"] will return ["root FileId", "foo FileId", "bar FileId"]
func traverseByPath(tx *gorm.DB, fileNames []string, user *User, fromHome bool) ([]File, error) {

	files := make([]File, len(fileNames)+1)

	var parentFolder *File
	var err error

	if fromHome {
		parentFolder, err = GetHomeFolder(tx, user)
	} else {
		parentFolder, err = GetRootFolder(tx)
	}

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
func CreateFolderByPath(tx *gorm.DB, path string, user *User, server string, fromHome bool) (*File, error) {

	paths := pathStringToArray(path, true)

	files, err := traverseByPath(tx, paths[0:(len(paths)-1)], user, fromHome)

	if err != nil {
		return nil, err
	}

	newFolder, err := CreateFolderByParentId(tx, files[len(files)-1].FileId, paths[len(paths)-1], user, server)

	if err != nil {
		return nil, err
	}

	return newFolder, nil
}

// CreateInitialFile if no error, will continue to call create Fragment for each and update
func CreateInitialFile(tx *gorm.DB, file *File, fileKey, fileIv, dataKey, dataIV string,
	user *User) error {

	// Verify that the file has the following properties
	/*
		1. FileName
		2. Size
		3. ParentFolderFileId
		4. Checksum

		If file.PasswordProtected, ensure that PasswordHint is there. Driver will hold the password.

		If file.ParityShards is not set (or lower than 0), use defaults
	*/

	if file.FileName == "" || file.Size == 0 || *file.ParentFolderFileId == "" || fileKey == "" || fileIv == "" || dataKey == "" || dataIV == "" {
		return ErrInvalidFile
	}

	// Verify that the user has permission to the folder they are writing to

	hasPermission, err := user.HasPermission(tx, &File{FileId: *file.ParentFolderFileId}, &PermissionNeeded{Read: true})
	if err != nil || !hasPermission {
		return err
	}

	hasPermission, err = user.HasPermission(tx, &File{FileId: *file.ParentFolderFileId}, &PermissionNeeded{Write: true})
	if err != nil || !hasPermission {
		return ErrNoPermission
	}

	// Check that there's no file with the same name in the same folder

	ls, err := ListFilesByFolderId(tx, *file.ParentFolderFileId, user)

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

	// Encrypt the data key and IV with the file key and IV

	newKey, err := EncryptWithKeyIV(dataKey, fileKey, fileIv)
	if err != nil {
		return err
	}
	newIv, err := EncryptWithKeyIV(dataIV, fileKey, fileIv)
	if err != nil {
		return err
	}

	newFile := &File{
		FileId:             file.FileId,
		FileName:           file.FileName,
		MIMEType:           mimetype,
		EntryType:          IsFile,
		ParentFolderFileId: file.ParentFolderFileId,
		VersionNo:          0,
		DataId:             uuid.New().String(),
		DataIdVersion:      0,
		Size:               file.Size,
		ActualSize:         file.Size,
		CreatedTime:        time.Now(),
		ModifiedUserUserId: &user.UserId,
		ModifiedTime:       time.Now(),
		VersioningMode:     file.VersioningMode,
		Checksum:           file.Checksum,
		TotalShards:        file.TotalShards,
		DataShards:         file.DataShards,
		ParityShards:       file.ParityShards,
		KeyThreshold:       file.KeyThreshold,
		PasswordProtected:  file.PasswordProtected,
		EncryptionKey:      newKey,
		EncryptionIv:       newIv,
		LastChecked:        time.Now(),
		Status:             FileStatusWriting,
		HandledServer:      file.HandledServer, // Should be the server Id of the server that is handling the file
	}

	// Create the file
	err = tx.Create(newFile).Error
	if err != nil {
		return err
	}

	*file = *newFile

	// Update and place into FileVersion
	return CreateFileVersionFromFile(tx, file, user)

}

// CreateFragment is called whenever the fragment has been written to disk
func CreateFragment(tx *gorm.DB, fileId string, dataId string, versionNo int, fragId int, serverId string, fragmentPath string) error {

	// Validating inputs
	if fileId == "" || dataId == "" || versionNo < 0 || fragId <= 0 || serverId == "" || fragmentPath == "" {
		return ErrInvalidAction
	}

	newFragment := Fragment{
		FileVersionFileId:    fileId,
		FileVersionDataId:    dataId,
		FileVersionVersionNo: versionNo,
		FragId:               fragId,
		ServerName:           serverId,
		FileFragmentPath:     fragmentPath,
		LastChecked:          time.Now(),
		Status:               FragmentStatusGood,
	}
	err := tx.Create(&newFragment).Error
	if err != nil {
		return err
	}

	return err
}

// FinishFile is called whenever a file has been all written (all fragments written)
func FinishFile(tx *gorm.DB, file *File, user *User, size, actualSize int64, checksum string) error {

	// Some checks
	if file.Status != FileStatusWriting {
		return ErrInvalidAction
	}
	if actualSize <= 0 {
		return ErrInvalidAction
	}

	// Update the file to be finished
	file.Status = FileStatusGood
	file.Checksum = checksum
	file.Size = size
	file.ActualSize = actualSize
	file.ModifiedUserUserId = &user.UserId
	file.ModifiedTime = time.Now()
	err := tx.Save(file).Error
	if err != nil {
		return err
	}

	// Updating the FileVersion
	return finaliseFileVersionFromFile(tx, file, size, actualSize)

}

// CreateFolderByParentId creates a folder based on the id given and returns the folder (File Object)
func CreateFolderByParentId(tx *gorm.DB, id string, folderName string, user *User, server string) (*File, error) {

	// Check that the id exists

	parentFolder := &File{FileId: id}

	err := tx.First(&parentFolder).Error

	if err != nil {
		return nil, ErrFileNotFound
	}

	return parentFolder.CreateSubFolder(tx, folderName, user, server)

}

// DeleteFileById deletes a file based on the FileId given.
func DeleteFileById(tx *gorm.DB, id string, user *User, server string) error {
	file, err := GetFileById(tx, id, user)
	if err != nil {
		return err
	}

	return file.Delete(tx, user, server)

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
func DeleteFolderByIdCascade(tx *gorm.DB, id string, user *User, server string) error {

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
					return deleteSubFoldersCascade(tx, &file, user, server)
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
func deleteSubFoldersCascade(tx *gorm.DB, file *File, user *User, server string) error {

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
				return deleteSubFoldersCascade(tx, &file, user, server)
			})
			if err != nil {
				return err
			}
		}
	}

	err := tx.Transaction(
		func(tx *gorm.DB) error {
			if err2 := deleteFileVersionFromFile(tx, file, server); err2 != nil {
				return err2
			}
			if err2 := deleteFilePermissions(tx, file); err2 != nil {
				return err2
			}
			if err2 := tx.Delete(file).Error; err2 != nil {
				return err2
			}
			return nil
		},
	)

	return err

}

// GetFileFragments returns the fragments associated with the File
func (f File) GetFileFragments(tx *gorm.DB, user *User) ([]Fragment, error) {

	// Check if the user has read file permissions to it

	if user != nil {
		hasPermission, err := user.HasPermission(tx, &f, &PermissionNeeded{Read: true})
		if err != nil {
			return nil, err
		} else if !hasPermission {
			return nil, ErrFileNotFound
		}
	} else {
		// Check if the file is public
		var count int64
		tx.Model(&SharedLink{}).Where("file_id = ?", f.FileId).Count(&count)
		if count == 0 {
			return nil, ErrFileNotFound
		}
	}

	var fragments []Fragment
	err := tx.Where("file_version_data_id = ?", f.DataId).Find(&fragments).Error
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
func (f *File) CreateSubFolder(tx *gorm.DB, folderName string, user *User, server string) (*File, error) {

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
		EntryType:          IsFolder,
		ParentFolder:       f,
		ParentFolderFileId: &f.FileId,
		VersionNo:          0,
		DataId:             "",
		DataIdVersion:      0,
		Size:               0,
		ActualSize:         0,
		CreatedTime:        time.Time{},
		ModifiedUserUserId: &user.UserId,
		ModifiedTime:       time.Time{},
		VersioningMode:     VersioningOff,
		Status:             1,
		HandledServer:      server,
	}

	// Transaction

	err = tx.Transaction(func(tx *gorm.DB) error {
		err = tx.Save(newFolder).Error

		if err != nil {
			return err
		}

		err = CreatePermissions(tx, newFolder)
		if err != nil {
			return err
		}

		return CreateFileVersionFromFile(tx, newFolder, user)

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

	changesMade := false

	if modificationsRequested.PasswordModification {

		pp, err := f.GetPasswordProtect(tx, user)
		if err != nil {
			return err
		}

		if modificationsRequested.NewPassword == "" && pp.PasswordActive {
			if f.PasswordUnprotect(tx, modificationsRequested.OldPassword, user) != nil {
				return ErrIncorrectPassword
			}
		} else if modificationsRequested.NewPassword == "" && !pp.PasswordActive {
			// No password change requested
		} else if modificationsRequested.NewPassword != "" {
			if f.PasswordProtect(tx, modificationsRequested.OldPassword,
				modificationsRequested.NewPassword, modificationsRequested.PasswordHint,
				user) != nil {
				return err
			}
		}

	} else {
		if modificationsRequested.FileName != "" {
			if f.FileName != modificationsRequested.FileName {
				err := f.rename(tx, modificationsRequested.FileName)
				if err != nil {
					return err
				}
				changesMade = true
			}
		}
		if modificationsRequested.MIMEType != "" {
			if f.MIMEType != modificationsRequested.MIMEType {
				changesMade = true
				f.MIMEType = modificationsRequested.MIMEType
			}
		}
		if modificationsRequested.VersioningMode >= VersioningOff &&
			modificationsRequested.VersioningMode <= VersioningOnDeltas {
			if f.VersioningMode != modificationsRequested.VersioningMode {
				changesMade = true
				f.VersioningMode = modificationsRequested.VersioningMode
			}
		}
	}

	if !changesMade {
		return nil
	}

	f.VersionNo = f.VersionNo + 1

	// increment version number and save it in FileVersion
	// Attempt to save and if it fails, rollback

	err = tx.Transaction(func(tx *gorm.DB) error {
		err = CreateFileVersionFromFile(tx, f, user)
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
	return tx.Save(f).Error
}

// PasswordProtect
// Adds a password to a file
func (f *File) PasswordProtect(tx *gorm.DB, oldPassword, newPassword, hint string, user *User) error {

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

	var passwordProtect PasswordProtect
	err = tx.Model(&PasswordProtect{}).Where("file_id = ?", f.FileId).Find(&passwordProtect).Error
	if err != nil {
		return err
	}

	err = passwordProtect.encryptWithPassword(tx, oldPassword, newPassword, hint)

	f.PasswordProtected = true
	return tx.Save(f).Error
}

// PasswordUnprotect
// Removes a password to a file
func (f *File) PasswordUnprotect(tx *gorm.DB, password string, user *User) error {

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

	var passwordProtect PasswordProtect
	err = tx.Model(&PasswordProtect{}).Where("file_id = ?", f.FileId).Find(&passwordProtect).Error
	if err != nil {
		return err
	}

	err = passwordProtect.removePassword(tx, password)

	f.PasswordProtected = false
	return tx.Save(f).Error
}

func (f *File) GetPasswordProtect(tx *gorm.DB, user *User) (*PasswordProtect, error) {

	// Check if the user has read permission (if not 404)
	hasPermissions, err := user.HasPermission(tx, f, &PermissionNeeded{Read: true})
	if err != nil {
		return nil, err
	} else if !hasPermissions {
		return nil, ErrFileNotFound
	}

	var passwordProtect PasswordProtect
	err = tx.Model(&PasswordProtect{}).Where("file_id = ?", f.FileId).Find(&passwordProtect).Error
	if err != nil {
		return nil, err
	}

	return &passwordProtect, nil
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

	// Check that the destination is a folder
	if newParent.EntryType != IsFolder {
		return ErrNotFolder
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
	f.ParentFolderFileId = &newParent.FileId
	f.VersionNo = f.VersionNo + 1

	err = tx.Transaction(func(tx *gorm.DB) error {
		err2 := tx.Save(f).Error
		if err2 != nil {
			return err2
		}

		err2 = CreatePermissions(tx, f)
		if err2 != nil {
			return err2
		}

		err2 = CreateFileVersionFromFile(tx, f, user)

		return err2

	})

	return err

}

// Copy copies the file to a new folder
func (f *File) Copy(tx *gorm.DB, newParent *File, user *User, server string) error {

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

	// Check that the user has write permissions on the new parent folder
	hasPermissions, err = user.HasPermission(tx, newParent, &PermissionNeeded{Write: true})
	if err != nil {
		return err
	} else if !hasPermissions {
		return ErrNoPermission
	}

	// Create a new file

	newFile := File{
		FileId:             uuid.New().String(),
		FileName:           f.FileName,
		MIMEType:           f.MIMEType,
		EntryType:          f.EntryType,
		ParentFolderFileId: &newParent.FileId,
		VersionNo:          0,
		DataId:             f.DataId,
		DataIdVersion:      f.DataIdVersion,
		Size:               f.Size,
		ActualSize:         f.ActualSize,
		CreatedTime:        time.Now(),
		ModifiedUserUserId: &user.UserId,
		ModifiedTime:       time.Now(),
		VersioningMode:     f.VersioningMode,
		Checksum:           f.Checksum,
		TotalShards:        f.TotalShards,
		DataShards:         f.DataShards,
		ParityShards:       f.ParityShards,
		KeyThreshold:       f.KeyThreshold,
		PasswordProtected:  f.PasswordProtected,
		EncryptionKey:      f.EncryptionKey,
		EncryptionIv:       f.EncryptionIv,
		LastChecked:        time.Now(),
		Status:             FileStatusGood,
		HandledServer:      server,
	}

	// Find the original passwordProtect and duplicate it
	var ogPP PasswordProtect
	err = tx.Model(&PasswordProtect{}).Where("file_id = ?", f.FileId).Find(&ogPP).Error
	if err != nil {
		return err
	}

	newPasswordProtect := PasswordProtect{
		FileId:         newFile.FileId,
		FileKey:        ogPP.FileKey,
		FileIv:         ogPP.FileIv,
		PasswordActive: ogPP.PasswordActive,
		PasswordSalt:   ogPP.PasswordSalt,
		PasswordNonce:  ogPP.PasswordNonce,
		PasswordHint:   ogPP.PasswordHint,
	}

	err = tx.Transaction(func(tx2 *gorm.DB) error {

		err2 := tx2.Save(&newFile).Error
		if err2 != nil {
			return err2
		}

		err2 = tx2.Save(&newPasswordProtect).Error
		if err2 != nil {
			return err2
		}

		err2 = CreatePermissions(tx2, &newFile)
		if err2 != nil {
			return err2
		}
		err2 = CreateFileVersionFromFile(tx2, &newFile, user)
		if err2 != nil {
			return err2
		}

		dc := DataCopies{
			DataId: newFile.DataId,
		}

		return tx2.Clauses(clause.OnConflict{DoNothing: true}).Save(&dc).Error
	})

	if err != nil {
		return err
	}

	// Recursive call if folder
	fileOrFolder, err := f.IsFileOrEmptyFolder(tx, user)
	if err != nil {
		return err
	}

	// Recursively go through and copy all items visible to user
	if !fileOrFolder {
		ls, err := f.ListContents(tx, user)
		if err != nil {
			return err
		}
		for _, item := range ls {
			err = item.Copy(tx, &newFile, user, server)
			if err != nil {
				return err
			}
		}
	}

	return err

}

// Delete deletes a file or folder and all of its contents
func (f *File) Delete(tx *gorm.DB, user *User, server string) error {

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
					return deleteSubFoldersCascade(tx, &file, user, server)
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
		err2 := deleteFileVersionFromFile(tx, f, server)
		if err2 != nil {
			return err2
		}
		err2 = deleteFilePermissions(tx, f)
		if err != nil {
			return err
		}
		// Delete Favorites and Shares
		err2 = tx.Where("file_id = ?", f.FileId).Delete(&FavoriteFileItems{}).Error
		if err2 != nil {
			return err2
		}
		err2 = tx.Where("file_id = ?", f.FileId).Delete(&SharedLink{}).Error
		if err2 != nil {
			return err2
		}
		err2 = tx.Where("file_id = ?", f.FileId).Delete(&SharedWithUser{}).Error
		if err2 != nil {
			return err2
		}
		err2 = tx.Where("file_id = ?", f.FileId).Delete(&SharedWithGroup{}).Error
		if err2 != nil {
			return err2
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
		return upsertUsersPermission(tx, f, permission, requestUser, users...)
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
		return upsertGroupsPermission(tx, f, permission, requestUser, groups...)
	} else {
		return ErrNoPermission
	}

}

// RemovePermission revokes permissions for that user or group.
func (f *File) RemovePermission(tx *gorm.DB, permission *Permission, user *User) error {

	var err error

	// Check if it's user or group
	if permission.UserId != nil {
		// User
		err = revokeUsersPermission(tx, f, []User{User{UserId: *permission.UserId}}, user)

	} else {
		// Group
		err = revokeGroupsPermission(tx, f, []Group{Group{GroupId: *permission.GroupId}})
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

	if oldPermission.UserId != nil {
		newUser, err := GetUserById(tx, *newPermission.UserId)
		if err != nil {
			return err
		}
		err = upsertUsersPermission(tx, f, PermissionNeeded, user, *newUser)
	} else {
		newGroup, err := GetGroupBasedOnGroupId(tx, *newPermission.GroupId)
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

	err = tx.Where("file_id = ?", f.FileId).Preload("User").Preload("Group").Find(&permissions).Error

	return permissions, err
}

// GetPermissionById returns a single permission for a file based on the permissionId
func (f *File) GetPermissionById(tx *gorm.DB, permissionId string, user *User) (*Permission, error) {
	// Check if user has read permission (if not 404)
	hasPermissions, err := user.HasPermission(tx, f, &PermissionNeeded{Read: true})
	if err != nil {
		return nil, err
	} else if !hasPermissions {
		return nil, ErrFileNotFound
	}

	var permission Permission

	err = tx.Where("file_id = ? AND permission_id = ?", f.FileId, permissionId).Find(&permission).Error
	if err != nil {
		return nil, err
	}

	return &permission, err

}

// UpdateFile
// At this stage, the server should have the
// updated file (already processed) , uncompressed file size, compressed data size, and key
func (f *File) UpdateFile(tx *gorm.DB, newSize int64, newActualSize int64, checksum string, handlingServer string, dataKey string, dataIv string, password string, user *User, newFileName string) error {

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

	// Get the PasswordProtect for the file
	var passwordProtect PasswordProtect
	err = tx.Model(&PasswordProtect{}).Where("file_id = ?", f.FileId).First(&passwordProtect).Error
	if err != nil {
		return err
	}

	fileKey := passwordProtect.FileKey
	fileIv := passwordProtect.FileIv

	if passwordProtect.PasswordActive && password != "" {
		fileKey, fileIv, err = passwordProtect.DecryptWithPassword(password)
		if err != nil {
			return err
		}
	} else if passwordProtect.PasswordActive && password == "" {
		return ErrPasswordRequired
	}

	// Encrypt the dataKey and dataIv
	dataKey, err = EncryptWithKeyIV(dataKey, fileKey, fileIv)
	if err != nil {
		return err
	}
	dataIv, err = EncryptWithKeyIV(dataIv, fileKey, fileIv)
	if err != nil {
		return err
	}

	if newFileName != "" {
		f.FileName = newFileName
	}

	// Update the file
	f.VersionNo = f.VersionNo + 1
	f.DataId = uuid.New().String()
	f.DataIdVersion = f.DataIdVersion + 1
	f.Size = newSize
	f.ActualSize = newActualSize
	f.ModifiedUserUserId = &user.UserId
	f.ModifiedTime = time.Now()
	f.Checksum = checksum
	f.LastChecked = time.Now()
	f.EncryptionKey = dataKey
	f.EncryptionIv = dataIv
	f.Status = FileStatusWriting
	f.HandledServer = handlingServer

	// Updating the file in the database as well as FileVersion

	err = tx.Transaction(func(tx *gorm.DB) error {
		err2 := tx.Save(f).Error
		if err2 != nil {
			return err2
		}
		// Create a new FileVersion
		err2 = CreateFileVersionFromFile(tx, f, user)
		return err2
	})
	return err
}

// UpdateFragment Called on each fragment to be created.
// TODO: Update this to call CreateFragment instead.
func (f *File) UpdateFragment(tx *gorm.DB, fragmentId int, fileFragmentPath string, checksum string, serverId string) error {
	frag := Fragment{
		FileVersionFileId:        f.FileId,
		FileVersionDataId:        f.DataId,
		FileVersionVersionNo:     f.VersionNo,
		FileVersionDataIdVersion: f.DataIdVersion,
		FragId:                   fragmentId,
		ServerName:               serverId,
		FileFragmentPath:         fileFragmentPath,
		LastChecked:              time.Now(),
		TotalShards:              f.TotalShards,
		Status:                   FragmentStatusGood,
	}

	return tx.Save(&frag).Error

}

// FinishUpdateFile
// Once all fragments are updated, the file is marked as finished.
// Updates Status to FileStatusGood, and LastChecked to now.
func (f *File) FinishUpdateFile(tx *gorm.DB, checksum string) error {

	err := tx.Transaction(func(tx *gorm.DB) error {
		f.Status = FileStatusGood
		f.LastChecked = time.Now()
		f.Checksum = checksum

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

		fileVersion.Status = FileStatusGood
		fileVersion.Size = f.Size
		fileVersion.ActualSize = f.ActualSize
		fileVersion.LastChecked = f.LastChecked

		return tx.Save(&fileVersion).Error

	})
	if err != nil {
		return err
	}

	if f.VersioningMode == VersioningOff {
		// Mark previous fragment as to be deleted.
		err = tx.Model(&FileVersion{}).Where("file_id = ? AND versioning_mode = ? AND data_id <> ?",
			f.FileId, VersioningOff, f.DataId).Update("status", FileStatusToBeDeleted).Error

	} else if f.VersioningMode == VersioningOnVersions {
		err = nil
	} else if f.VersioningMode == VersioningOnDeltas {
		panic("Not implemented")
	}

	return err

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
	if f.EntryType != IsFolder {
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
	if f.EntryType == IsFile {
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

// GetDecryptionKey returns the Key and IV of a file given a password (or not)
func (f *File) GetDecryptionKey(tx *gorm.DB, user *User, password string) (string, string, error) {

	if user != nil {
		// Check if user has read permission (if not 404)
		hasPermissions, err := user.HasPermission(tx, f, &PermissionNeeded{Read: true})
		if !hasPermissions {
			return "", "", ErrFileNotFound
		} else if err != nil {
			return "", "", err
		}
	} else {
		// Check if the file is public
		var count int64
		tx.Model(&SharedLink{}).Where("file_id = ?", f.FileId).Count(&count)
		if count == 0 {
			return "", "", ErrFileNotFound
		}
	}

	// Get FileKey, FileIv from PasswordProtect

	var passwordProtect PasswordProtect
	err := tx.Model(&PasswordProtect{}).Where("file_id = ?", f.FileId).First(&passwordProtect).Error
	if err != nil {
		return "", "", err
	}

	// if password nil
	if password == "" && !passwordProtect.PasswordActive {

		// decrypt file key with PasswordProtect
		fileKey, err := DecryptWithKeyIV(f.EncryptionKey, passwordProtect.FileKey, passwordProtect.FileIv)
		if err != nil {
			return "", "", err
		}
		fileIv, err := DecryptWithKeyIV(f.EncryptionIv, passwordProtect.FileKey, passwordProtect.FileIv)
		if err != nil {
			return "", "", err
		}
		return fileKey, fileIv, nil
	} else if password == "" && passwordProtect.PasswordActive {
		return "", "", ErrPasswordRequired
	} else {
		decryptedFileKey, decryptedFileIv, err := passwordProtect.DecryptWithPassword(password)
		if err != nil {
			return "", "", err
		}
		fileKey, err := DecryptWithKeyIV(f.EncryptionKey, decryptedFileKey, decryptedFileIv)
		if err != nil {
			return "", "", err
		}
		fileIv, err := DecryptWithKeyIV(f.EncryptionIv, decryptedFileKey, decryptedFileIv)
		if err != nil {
			return "", "", err
		}
		return fileKey, fileIv, nil

	}

}

// GetAllVersions returns all versions of a file
func (f *File) GetAllVersions(tx *gorm.DB, user *User) ([]FileVersion, error) {

	// Check if user has read permission (if not 404)
	hasPermissions, err := user.HasPermission(tx, f, &PermissionNeeded{Read: true})
	if !hasPermissions {
		return nil, ErrFileNotFound
	} else if err != nil {
		return nil, err
	}

	var fileVersions []FileVersion

	err = tx.Model(&FileVersion{}).Where("file_id = ? AND status <> ?", f.FileId, FileStatusToBeDeleted).Find(&fileVersions).Error

	if err != nil {
		return nil, err
	}

	return fileVersions, nil

}

// DeleteFileVersion marks a FileVersion for deletion.
func (f *File) DeleteFileVersion(tx *gorm.DB, user *User, versionNo int) error {

	// Check if user has read permission (if not 404)
	hasPermissions, err := user.HasPermission(tx, f, &PermissionNeeded{Read: true})
	if !hasPermissions {
		return ErrFileNotFound
	} else if err != nil {
		return err
	}

	// Check if user has write permission (if not 403)
	hasPermissions, err = user.HasPermission(tx, f, &PermissionNeeded{Write: true})
	if err != nil {
		return err
	} else if !hasPermissions {
		return ErrNoPermission
	}

	// Check if the version exists
	var fileVersion FileVersion
	err = tx.Model(&FileVersion{}).Where("file_id = ? AND version_no = ?", f.FileId, versionNo).First(&fileVersion).Error
	if err != nil {
		return ErrVersionNotFound
	}

	// Mark version for deletion
	return tx.Model(&FileVersion{}).Where("file_id = ? AND version_no = ?",
		f.FileId, versionNo).Update("status", FileStatusToBeDeleted).Error

}

// GetPath returns the path of a file with fileID
func (f *File) GetPath(tx *gorm.DB, user *User) ([]File, error) {

	// Check if user has read permission (if not 404)
	hasPermissions, err := user.HasPermission(tx, f, &PermissionNeeded{Read: true})
	if !hasPermissions {
		return nil, ErrFileNotFound
	} else if err != nil {
		return nil, err
	}

	files := []File{
		*f,
	}
	parentFile := f
	continueSearching := true
	for continueSearching {
		pfid := *parentFile.ParentFolderFileId
		var tempParent File
		if pfid == "00000000-0000-0000-0000-000000000000" {
			// at root
			err = tx.Where(&File{FileId: pfid}).First(&tempParent).Error
			files = append(files, tempParent)
			continueSearching = false
			break
		}

		err = tx.Where(&File{FileId: pfid}).First(&tempParent).Error
		if err != nil {
			return nil, err
		}
		hasPermission, err := user.HasPermission(tx, &tempParent, &PermissionNeeded{Read: true})
		if err != nil {
			return nil, err
		}
		if !hasPermission {
			continueSearching = false
			break
		}
		files = append(files, tempParent)
		parentFile = &tempParent
	}

	return files, nil
}

// CreateSharedLink creates a shared link for a file
// If link is not provided, it will generate a random one
func (f *File) CreateSharedLink(tx *gorm.DB, user *User, link string) (*SharedLink, error) {

	// Check if user has shared link permission
	permission, err := user.HasPermission(tx, f, &PermissionNeeded{Read: true, Share: true})
	if err != nil {
		return nil, err
	}
	if !permission {
		return nil, ErrNoPermission
	}

	if link == "" {
		// Generate a random link

		uniqueLink := false

		for !uniqueLink {
			rpg := random_phrase_generator.New()
			link = rpg.GenerateRandomPhrase()

			// Check that the link is unique
			var count int64
			tx.Model(&SharedLink{}).Where("shortened_link = ?", link).Count(&count)
			if count == 0 {
				uniqueLink = true
			}
		}

	} else {

		// Check that the new link passes REGEX test
		match, _ := regexp.MatchString(SharedURLRegexConstriant, link)
		if !match {
			return nil, ErrLinkConstraint
		}

		// Check that the link is unique
		var count int64
		tx.Model(&SharedLink{}).Where("shortened_link = ?", link).Count(&count)
		if count != 0 {
			return nil, ErrLinkExists
		}
	}

	// Create the shared link
	sharedLink := SharedLink{
		FileId:        f.FileId,
		ShortenedLink: link,
		CreatedTime:   time.Now(),
	}
	if tx.Create(&sharedLink).Error != nil {
		return nil, err
	}

	return &sharedLink, nil
}

// GetSharedLinks returns all shared links for a file
func (f *File) GetSharedLinks(tx *gorm.DB, user *User) ([]SharedLink, error) {
	// Check if user has read permission
	hasPermissions, err := user.HasPermission(tx, f, &PermissionNeeded{Read: true})
	if err != nil {
		return nil, err
	}
	if !hasPermissions {
		return nil, ErrFileNotFound
	}

	// Get the shared links
	var sharedLinks []SharedLink
	err = tx.Model(&SharedLink{}).Where("file_id = ?", f.FileId).Find(&sharedLinks).Error

	return sharedLinks, err
}

// DeleteSharedLink deletes a shared link for a file given it's link
func (f *File) DeleteSharedLink(tx *gorm.DB, user *User, link string) error {
	// Check if user has read permission
	hasPermissions, err := user.HasPermission(tx, f, &PermissionNeeded{Read: true})
	if err != nil {
		return err
	}
	if !hasPermissions {
		return ErrFileNotFound
	}

	// Check if user has shared link permission
	permission, err := user.HasPermission(tx, f, &PermissionNeeded{Read: true, Share: true})
	if err != nil {
		return err
	}
	if !permission {
		return ErrNoPermission
	}
	// Delete the shared link
	return tx.Where("file_id = ? AND shortened_link = ?", f.FileId, link).Delete(&SharedLink{}).Error
}

// UpdateSharedLink updates a shared link for a file with a new link
func (f *File) UpdateSharedLink(tx *gorm.DB, user *User, link string, newLink string) error {
	// Check if user has read permission
	hasPermissions, err := user.HasPermission(tx, f, &PermissionNeeded{Read: true})
	if err != nil {
		return err
	}
	if !hasPermissions {
		return ErrFileNotFound
	}

	// Check if user has shared link permission
	permission, err := user.HasPermission(tx, f, &PermissionNeeded{Read: true, Share: true})
	if err != nil {
		return err
	}
	if !permission {
		return ErrNoPermission
	}

	// Check that the new link passes REGEX test
	match, _ := regexp.MatchString(SharedURLRegexConstriant, newLink)
	if !match {
		return ErrLinkConstraint
	}

	// Check that the new link is unique

	var count int64
	err = tx.Model(&SharedLink{}).Where("shortened_link = ?", newLink).Count(&count).Error
	if err != nil {
		return err
	}
	if count != 0 {
		return ErrLinkExists
	}

	// Update the shared link
	return tx.Model(&SharedLink{}).Where("file_id = ? AND shortened_link = ?", f.FileId, link).
		Update("shortened_link", newLink).Error
}

// AddToFavorites will add the file to a user's favorites
func (f *File) AddToFavorites(tx *gorm.DB, user *User) error {

	// Check if the user has permisisons to read the file

	perm, err := user.HasPermission(tx, f, &PermissionNeeded{Read: true})
	if err != nil {
		return err
	} else if !perm {
		return ErrFileNotFound
	}

	return tx.Clauses(clause.OnConflict{DoNothing: true}).
		Save(&FavoriteFileItems{
			FileId: f.FileId,
			UserId: user.UserId,
		}).Error
}

// RemoveFromFavorites will remove the file from a user's favorites
func (f *File) RemoveFromFavorites(tx *gorm.DB, user *User) error {

	if err := tx.Model(&FavoriteFileItems{}).
		Where("file_id = ? AND user_id = ?", f.FileId, user.UserId).
		Delete(FavoriteFileItems{}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}

	return nil

}

// RotateKey will rotate the key and IV of a file
func (f *File) RotateKey(tx *gorm.DB, user *User, password string) error {

	// First, we need to extract out each file versions key and iv.

	var fileVersions []FileVersion
	if err := tx.Model(&FileVersion{}).Where("file_id = ?", f.FileId).Find(&fileVersions).Error; err != nil {
		return err
	}

	type KeyIvPair struct {
		Key string
		Iv  string
	}

	newKeyIvPair := make([]KeyIvPair, len(fileVersions))

	// Generate a new key and iv.

	newFileKey, newFileIv, err := GenerateKeyIV()
	if err != nil {
		return err
	}

	for i, fileVersion := range fileVersions {

		oldKey, oldIv, err := fileVersion.GetDecryptionKey(tx, user, password)
		if err != nil {
			return err
		}

		// encrypt it with the new key and iv
		newKey, err := EncryptWithKeyIV(oldKey, newFileKey, newFileIv)
		if err != nil {
			return err
		}
		newIv, err := EncryptWithKeyIV(oldIv, newFileKey, newFileIv)
		if err != nil {
			return err
		}

		newKeyIvPair[i].Key = newKey
		newKeyIvPair[i].Iv = newIv
	}

	latestVersionOldKey, latestVersionOldIv, err := f.GetDecryptionKey(tx, user, password)
	if err != nil {
		return err
	}

	latestVersionNewKey, err := EncryptWithKeyIV(latestVersionOldKey, newFileKey, newFileIv)
	if err != nil {
		return err
	}
	latestVersionNewIv, err := EncryptWithKeyIV(latestVersionOldIv, newFileKey, newFileIv)
	if err != nil {
		return err
	}

	// Encrypt the file key and iv if the password is not empty
	var nonce, salt []byte

	if password != "" {
		// Generate a secure salt
		salt = make([]byte, 32)
		_, err = rand.Read(salt)
		if err != nil {
			return err
		}

		key, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, 32)
		if err != nil {
			return err
		}

		// ENCRYPT
		fileKeyBytes, err := hex.DecodeString(newFileKey)
		if err != nil {
			return err
		}

		fileIvBytes, err := hex.DecodeString(newFileIv)
		if err != nil {
			return err
		}

		block, err := aes.NewCipher(key)
		if err != nil {
			return err
		}

		aesgcm, err := cipher.NewGCM(block)
		if err != nil {
			return err
		}

		nonce = make([]byte, aesgcm.NonceSize())
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			return err
		}

		var p PasswordProtect
		err = tx.Model(&p).Where("file_id = ?", f.FileId).First(&p).Error
		if err != nil {
			return err
		}

		p.PasswordNonce = hex.EncodeToString(nonce)

		newFileKeyBytes := aesgcm.Seal(nil, nonce, fileKeyBytes, nil)
		newFileIvBytes := aesgcm.Seal(nil, nonce, fileIvBytes, nil)

		newFileKey = hex.EncodeToString(newFileKeyBytes)
		newFileIv = hex.EncodeToString(newFileIvBytes)
	}

	// Update the db

	err = tx.Transaction(func(tx *gorm.DB) error {

		// Update file (latest version)
		err = tx.Model(&f).Where("file_id = ?", f.FileId).
			Update("encryption_key", latestVersionNewKey).
			Update("encryption_iv", latestVersionNewIv).Error
		if err != nil {
			return err
		}

		// Updating file versions
		for i, fileVersion := range fileVersions {
			if err := tx.Model(&fileVersion).Where("file_id = ? AND version_no = ?", f.FileId, fileVersion.VersionNo).
				Update("encryption_key", newKeyIvPair[i].Key).
				Update("encryption_iv", newKeyIvPair[i].Iv).Error; err != nil {
				return err
			}
		}

		// Update Password Protect
		if err := tx.Model(&PasswordProtect{}).Where("file_id = ?", f.FileId).
			Update("file_key", newFileKey).
			Update("file_iv", newFileIv).
			Update("password_nonce", hex.EncodeToString(nonce)).
			Update("password_salt", hex.EncodeToString(salt)).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}
