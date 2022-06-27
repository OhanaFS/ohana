package dbfs

import (
	"gorm.io/gorm"
	"time"
)

type FileVersion struct {
	FileID                string `gorm:"primaryKey"`
	VersionNo             uint   `gorm:"primaryKey"`
	FileName              string `gorm:"not null"`
	MIMEType              string
	EntryType             int8         `gorm:"not null"`
	ParentFolder          *FileVersion `gorm:"foreignKey:ParentFolderFileID,ParentFolderVersionNo"`
	ParentFolderFileID    string
	ParentFolderVersionNo uint
	DataID                string
	DataIDVersion         uint
	Size                  uint      `gorm:"not null"`
	ActualSize            uint      `gorm:"not null"`
	CreatedTime           time.Time `gorm:"not null"`
	ModifiedUser          User      `gorm:"foreignKey:ModifiedUserUserID"`
	ModifiedUserUserID    string
	ModifiedTime          time.Time `gorm:"not null; autoUpdateTime"`
	VersioningMode        int8      `gorm:"not null"`
	Checksum              string
	FragCount             uint8
	ParityCount           uint8
	EncryptionKey         string
	PasswordProtected     bool
	PasswordHint          string
	LinkFile              *FileVersion `gorm:"foreignKey:LinkFileFileID,LinkFileVersionNo"`
	LinkFileFileID        string
	LinkFileVersionNo     uint
	LastChecked           time.Time
	Status                int8   `gorm:"not null"`
	HandledServer         string `gorm:"not null"`
	Patch                 bool
	PatchBaseVersion      uint
}

// createFileVersionFromFile creates a FileVersion from a File
func createFileVersionFromFile(tx *gorm.DB, file *File, user *User) error {

	// Get the current parent folder and it's current version

	parentFolder, err := GetFileByID(tx, file.ParentFolderFileID, user)
	if err != nil {
		return err
	}

	fileVersion := FileVersion{
		FileID:                file.FileID,
		VersionNo:             file.VersionNo,
		FileName:              file.FileName,
		MIMEType:              file.MIMEType,
		EntryType:             file.EntryType,
		ParentFolderFileID:    file.ParentFolderFileID,
		ParentFolderVersionNo: parentFolder.VersionNo,
		DataID:                file.DataID,
		DataIDVersion:         file.DataIDVersion,
		Size:                  file.Size,
		ActualSize:            file.ActualSize,
		CreatedTime:           file.CreatedTime,
		ModifiedUserUserID:    file.ModifiedUserUserID,
		ModifiedTime:          file.ModifiedTime,
		VersioningMode:        file.VersioningMode,
		Checksum:              file.Checksum,
		FragCount:             file.FragCount,
		ParityCount:           file.ParityCount,
		EncryptionKey:         file.EncryptionKey,
		PasswordProtected:     file.PasswordProtected,
		PasswordHint:          file.PasswordHint,
		//LinkFileFileID:        "GET LINKED FOLDER", // NOT READY
		//LinkFileVersionNo:     0,                   // NOT READY
		LastChecked:   file.LastChecked,
		Status:        file.Status,
		HandledServer: file.HandledServer,
	}

	return tx.Save(&fileVersion).Error

}

// finaliseFileVersionFromFile finalises the status to be done (FILESTATUSGOOD)
func finaliseFileVersionFromFile(tx *gorm.DB, file *File) error {
	return tx.Model(&FileVersion{}).Where("file_id = ? AND version_no = ?", file.FileID, file.VersionNo).
		Update("status", FILESTATUSGOOD).Error
}

// getFileVersionFromFile returns the version requested of a file/
func getFileVersionFromFile(tx *gorm.DB, file *File, version int) (*FileVersion, error) {
	var fileVersion *FileVersion
	err := tx.Where("file_id = ? AND version_no = ?", file.FileID, version).First(fileVersion).Error
	if err != nil {
		return nil, err
	}
	return fileVersion, nil
}

// GetFragments returns the fragments of a FileVersion
func (fileVersion *FileVersion) GetFragments(tx *gorm.DB, user *User) ([]Fragment, error) {

	// Check if user has permissions to the file

	// We need to check if the user has permission to the current file.
	// This actually might be an issue with deleted files. Need to think how to solve this issue.
	// Possibly changing HasPermission() to check from FileVersion instead of File. Thus it should have records
	// of deleted files as well. Or for now, we just take it that you can't see deleted files.

	// Get original File and check permissions (can't get file without permissions)
	file, err := GetFileByID(tx, fileVersion.FileID, user)
	if err != nil {
		return nil, err
	}

	// Check if the file is a folder

	if file.EntryType == ISFOLDER {
		return nil, ErrNotFile
	}

	var fragments []Fragment
	err = tx.Model(&fragments).Where("file_id = ? AND version_no = ?", fileVersion.FileID, fileVersion.VersionNo).Find(&fragments).Error
	if err != nil {
		return nil, err
	}
	return fragments, nil
}

// ListFiles list the files in a folder of a FileVersion
func (fileVersion *FileVersion) ListFiles(tx *gorm.DB, user *User) ([]FileVersion, error) {

	// Check if user has permissions to the file

	// We need to check if the user has permission to the current file.
	// This actually might be an issue with deleted files. Need to think how to solve this issue.
	// Possibly changing HasPermission() to check from FileVersion instead of File. Thus it should have records
	// of deleted files as well. Or for now, we just take it that you can't see deleted files.

	// Get original File and check permissions (can't get file without permissions)
	file, err := GetFileByID(tx, fileVersion.FileID, user)
	if err != nil {
		return nil, err
	}

	// Check if the file is a folder

	if file.EntryType == ISFILE {
		return nil, ErrNotFolder
	}

	var files []FileVersion
	err = tx.Model(&files).Where("parent_folder_file_id = ? AND parent_folder_version_no = ?", fileVersion.FileID, fileVersion.VersionNo).Find(&files).Error
	if err != nil {
		return nil, err
	}
	return files, nil
}

//deleteFileVersionFromFile will mark all versions as to be deleted, but will not delete
// till the system clears it up as a chron job.
func deleteFileVersionFromFile(tx *gorm.DB, file *File) error {

	// First, we'll create a new history entry to show when the file was deleted with timestamp

	// Get the current parent folder and it's current version
	parentFolder, err := GetFileByID(tx, file.ParentFolderFileID, nil)
	if err != nil {
		return err
	}

	fileVersion := FileVersion{
		FileID:                file.FileID,
		VersionNo:             file.VersionNo + 1,
		FileName:              file.FileName,
		MIMEType:              file.MIMEType,
		EntryType:             file.EntryType,
		ParentFolderFileID:    file.ParentFolderFileID,
		ParentFolderVersionNo: parentFolder.VersionNo,
		DataID:                file.DataID,
		DataIDVersion:         file.DataIDVersion,
		Size:                  file.Size,
		ActualSize:            file.ActualSize,
		CreatedTime:           file.CreatedTime,
		ModifiedUserUserID:    file.ModifiedUserUserID,
		ModifiedTime:          file.ModifiedTime,
		VersioningMode:        file.VersioningMode,
		Checksum:              file.Checksum,
		FragCount:             file.FragCount,
		ParityCount:           file.ParityCount,
		EncryptionKey:         file.EncryptionKey,
		PasswordProtected:     file.PasswordProtected,
		PasswordHint:          file.PasswordHint,
		//LinkFileFileID:        "GET LINKED FOLDER", // NOT READY
		//LinkFileVersionNo:     0,                   // NOT READY
		LastChecked:   file.LastChecked,
		Status:        FILESTATUSDELETED,
		HandledServer: file.HandledServer,
	}

	err = tx.Save(&fileVersion).Error
	if err != nil {
		return err
	}

	// Next, we'll mark everything as deleted.

	err = tx.Model(&FileVersion{}).Where("file_id = ?", file.FileID).Update("status", FILESTATUSDELETED).Error
	if err != nil {
		return err
	}

	return nil
}
