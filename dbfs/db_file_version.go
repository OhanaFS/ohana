package dbfs

import (
	"gorm.io/gorm"
	"time"
)

type FileVersion struct {
	FileId                string       `gorm:"primaryKey" json:"file_id"`
	VersionNo             int          `gorm:"primaryKey" json:"version_no"`
	FileName              string       `gorm:"not null" json:"file_name"`
	MIMEType              string       `json:"mime_type"`
	EntryType             int8         `gorm:"not null" json:"entry_type"`
	ParentFolder          *FileVersion `gorm:"foreignKey:ParentFolderFileId,ParentFolderVersionNo" json:"'-'"`
	ParentFolderFileId    *string      `json:"parent_folder_id"`
	ParentFolderVersionNo *int         `json:"-"`
	DataId                string       `json:"-"`
	DataIdVersion         int          `json:"data_version_no"`
	Size                  int          `gorm:"not null" json:"size"`
	ActualSize            int          `gorm:"not null" json:"actual_size"`
	CreatedTime           time.Time    `gorm:"not null" json:"created_time"`
	ModifiedUser          User         `gorm:"foreignKey:ModifiedUserUserId" json:"-"`
	ModifiedUserUserId    *string      `json:"modified_user_user_id"`
	ModifiedTime          time.Time    `gorm:"not null; autoUpdateTime" json:"modified_time"`
	VersioningMode        int8         `gorm:"not null" json:"versioning_mode"`
	Checksum              string       `json:"checksum"`
	TotalShards           int          `json:"total_shards"`
	DataShards            int          `json:"data_shards"`
	ParityShards          int          `json:"parity_shards"`
	KeyThreshold          int          `json:"key_threshold"`
	EncryptionKey         string       `json:"-"`
	EncryptionIv          string       `json:"-"`
	PasswordProtected     bool         `json:"password_protected"`
	LinkFile              *FileVersion `gorm:"foreignKey:LinkFileFileId,LinkFileVersionNo" json:"-"`
	LinkFileFileId        *string      `json:"link_file_id"`
	LinkFileVersionNo     *int         `json:"-"`
	LastChecked           time.Time    `json:"last_checked"`
	Status                int8         `gorm:"not null" json:"status"`
	HandledServer         string       `gorm:"not null" json:"-"`
	Patch                 bool         `json:"-"`
	PatchBaseVersion      int          `json:"-"`
}

// CreateFileVersionFromFile creates a FileVersion from a File
func CreateFileVersionFromFile(tx *gorm.DB, file *File, user *User) error {

	// Get the current parent folder and it's current version

	parentFolder, err := GetFileById(tx, *file.ParentFolderFileId, user)
	if err != nil {
		return err
	}

	fileVersion := FileVersion{
		FileId:                file.FileId,
		VersionNo:             file.VersionNo,
		FileName:              file.FileName,
		MIMEType:              file.MIMEType,
		EntryType:             file.EntryType,
		ParentFolderFileId:    file.ParentFolderFileId,
		ParentFolderVersionNo: &parentFolder.VersionNo,
		DataId:                file.DataId,
		DataIdVersion:         file.DataIdVersion,
		Size:                  file.Size,
		ActualSize:            file.ActualSize,
		CreatedTime:           file.CreatedTime,
		ModifiedUserUserId:    file.ModifiedUserUserId,
		ModifiedTime:          file.ModifiedTime,
		VersioningMode:        file.VersioningMode,
		Checksum:              file.Checksum,
		TotalShards:           file.TotalShards,
		DataShards:            file.DataShards,
		ParityShards:          file.ParityShards,
		KeyThreshold:          file.KeyThreshold,
		EncryptionKey:         file.EncryptionKey,
		EncryptionIv:          file.EncryptionIv,
		PasswordProtected:     file.PasswordProtected,
		//LinkFileFileId:        "GET LINKED FOLDER", // NOT READY
		//LinkFileVersionNo:     0,                   // NOT READY
		LastChecked:   file.LastChecked,
		Status:        file.Status,
		HandledServer: file.HandledServer,
	}

	return tx.Save(&fileVersion).Error

}

// finaliseFileVersionFromFile finalises the status to be done (FileStatusGood)
func finaliseFileVersionFromFile(tx *gorm.DB, file *File) error {
	return tx.Model(&FileVersion{}).Where("file_id = ? AND version_no = ?", file.FileId, file.VersionNo).
		Update("status", FileStatusGood).Error
}

// getFileVersionFromFile returns the version requested of a file/
func getFileVersionFromFile(tx *gorm.DB, file *File, version int) (*FileVersion, error) {
	var fileVersion FileVersion
	err := tx.Where("file_id = ? AND version_no = ?", file.FileId, version).First(&fileVersion).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrVersionNotFound
		}
		return nil, err
	}
	return &fileVersion, nil
}

// GetFragments returns the fragments of a FileVersion
func (fv *FileVersion) GetFragments(tx *gorm.DB, user *User) ([]Fragment, error) {

	// Check if user has permissions to the file

	// We need to check if the user has permission to the current file.
	// This actually might be an issue with deleted files. Need to think how to solve this issue.
	// Possibly changing HasPermission() to check from FileVersion instead of File. Thus it should have records
	// of deleted files as well. Or for now, we just take it that you can't see deleted files.

	// Get original File and check permissions (can't get file without permissions)
	file, err := GetFileById(tx, fv.FileId, user)
	if err != nil {
		return nil, err
	}

	// Check if the file is a folder

	if file.EntryType == IsFolder {
		return nil, ErrNotFile
	}

	var fragments []Fragment
	err = tx.Model(&fragments).Where("file_version_file_id = ? AND file_version_data_id = ?", fv.FileId, fv.DataId).Find(&fragments).Error
	if err != nil {
		return nil, err
	}
	return fragments, nil
}

// ListFiles list the files in a folder of a FileVersion
func (fv *FileVersion) ListFiles(tx *gorm.DB, user *User) ([]FileVersion, error) {

	// Check if user has permissions to the file

	// We need to check if the user has permission to the current file.
	// This actually might be an issue with deleted files. Need to think how to solve this issue.
	// Possibly changing HasPermission() to check from FileVersion instead of File. Thus it should have records
	// of deleted files as well. Or for now, we just take it that you can't see deleted files.

	// Get original File and check permissions (can't get file without permissions)
	file, err := GetFileById(tx, fv.FileId, user)
	if err != nil {
		return nil, err
	}

	// Check if the file is a folder

	if file.EntryType == IsFile {
		return nil, ErrNotFolder
	}

	var files []FileVersion
	err = tx.Model(&files).Where("parent_folder_file_id = ? AND parent_folder_version_no = ?", fv.FileId, fv.VersionNo).Find(&files).Error
	if err != nil {
		return nil, err
	}
	return files, nil
}

//deleteFileVersionFromFile will mark all versions as to be deleted, but will not delete
// till the system clears it up as a chron job.
func deleteFileVersionFromFile(tx *gorm.DB, file *File, server string) error {

	// First, we'll create a new history entry to show when the file was deleted with timestamp

	// Get the current parent folder and it's current version
	parentFolder := File{FileId: *file.ParentFolderFileId}
	err := tx.First(&parentFolder).Error
	if err != nil {
		return err
	}

	var status int8

	if file.EntryType == IsFolder {
		status = FileStatusDeleted
	} else {
		status = FileStatusToBeDeleted
	}

	fileVersion := FileVersion{
		FileId:                file.FileId,
		VersionNo:             file.VersionNo + 1,
		FileName:              file.FileName,
		MIMEType:              file.MIMEType,
		EntryType:             file.EntryType,
		ParentFolderFileId:    file.ParentFolderFileId,
		ParentFolderVersionNo: &parentFolder.VersionNo,
		DataId:                file.DataId,
		DataIdVersion:         file.DataIdVersion,
		Size:                  file.Size,
		ActualSize:            file.ActualSize,
		CreatedTime:           file.CreatedTime,
		ModifiedUserUserId:    file.ModifiedUserUserId,
		ModifiedTime:          time.Now(),
		VersioningMode:        file.VersioningMode,
		Checksum:              file.Checksum,
		TotalShards:           file.TotalShards,
		DataShards:            file.DataShards,
		ParityShards:          file.ParityShards,
		KeyThreshold:          file.KeyThreshold,
		EncryptionKey:         file.EncryptionKey,
		PasswordProtected:     file.PasswordProtected,
		//LinkFileFileId:        "GET LINKED FOLDER", // NOT READY
		//LinkFileVersionNo:     0,                   // NOT READY
		LastChecked:   file.LastChecked,
		Status:        status,
		HandledServer: server,
	}

	err = tx.Save(&fileVersion).Error
	if err != nil {
		return err
	}

	// Next, we'll mark everything as deleted.

	err = tx.Model(&FileVersion{}).Where("file_id = ?", file.FileId).Update("status", FileStatusToBeDeleted).Error
	if err != nil {
		return err
	}

	return nil
}

// GetDecryptionKey returns the Key and IV of a file given a password (or not)
func (fv *FileVersion) GetDecryptionKey(tx *gorm.DB, user *User, password string) (string, string, error) {

	// Check permissions for the original file
	_, err := GetFileById(tx, fv.FileId, user)
	if err != nil {
		return "", "", err
	}

	// Get FileKey, FileIv from PasswordProtect

	var passwordProtect PasswordProtect
	err = tx.Model(&PasswordProtect{}).Where("file_id = ?", fv.FileId).First(&passwordProtect).Error
	if err != nil {
		return "", "", err
	}

	// if password nil
	if password == "" && !passwordProtect.PasswordActive {

		// decrypt file key with PasswordProtect
		fileKey, err := DecryptWithKeyIV(fv.EncryptionKey, passwordProtect.FileKey, passwordProtect.FileIv)
		if err != nil {
			return "", "", err
		}
		fileIv, err := DecryptWithKeyIV(fv.EncryptionIv, passwordProtect.FileKey, passwordProtect.FileIv)
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
		fileKey, err := DecryptWithKeyIV(fv.EncryptionKey, decryptedFileKey, decryptedFileIv)
		if err != nil {
			return "", "", err
		}
		fileIv, err := DecryptWithKeyIV(fv.EncryptionIv, decryptedFileKey, decryptedFileIv)
		if err != nil {
			return "", "", err
		}
		return fileKey, fileIv, nil

	}

}

func (f *File) RevertFileToVersion(tx *gorm.DB, versionNo int, user *User) error {

	// Check permissions for the original file
	perm, err := user.HasPermission(tx, f, &PermissionNeeded{Write: true})
	if err != nil || !perm {
		return err
	}

	// Getting old version

	var oldVersion FileVersion

	err = tx.Model(&FileVersion{}).Where("file_id = ? AND version_no = ?", f.FileId, versionNo).First(&oldVersion).Error
	if err != nil {
		return err
	}

	// Setting it as the new version
	f.FileName = oldVersion.FileName
	f.MIMEType = oldVersion.MIMEType
	// f.ParentFolderFileID will not be updated.
	f.VersionNo = f.VersionNo + 1
	f.DataId = oldVersion.DataId
	f.DataIdVersion = f.DataIdVersion + 1
	f.Size = oldVersion.Size
	f.ActualSize = oldVersion.ActualSize
	// no need to update CreatedTime
	f.ModifiedUserUserId = &user.UserId
	f.ModifiedTime = time.Now()
	f.Checksum = oldVersion.Checksum
	f.TotalShards = oldVersion.TotalShards
	f.DataShards = oldVersion.DataShards
	f.ParityShards = oldVersion.ParityShards
	f.KeyThreshold = oldVersion.KeyThreshold
	f.EncryptionKey = oldVersion.EncryptionKey
	f.EncryptionIv = oldVersion.EncryptionIv
	f.LastChecked = oldVersion.LastChecked
	f.Status = oldVersion.Status

	// Save
	err = tx.Save(f).Error
	if err != nil {
		return err
	}
	// Update the file
	err = CreateFileVersionFromFile(tx, f, user)

	return err

}
