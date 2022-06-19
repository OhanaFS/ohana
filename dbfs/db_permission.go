package dbfs

import (
	"gorm.io/gorm"
	"time"
)

type Permission struct {
	FileID       string `gorm:"primaryKey"`
	PermissionID uint   `gorm:"primaryKey;autoIncrement"`
	User         User   `gorm:"foreignKey:UserID"`
	UserID       string
	Group        Group `gorm:"foreignKey:GroupID"`
	GroupID      string
	CanRead      bool
	CanWrite     bool
	CanExecute   bool
	CanShare     bool
	VersionNo    uint
	Audit        bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Status       int
}

type PermissionInterface interface {

	// Global Private Functions. Called only by File

	// createPermissions takes in the parent file and copies it to the new File
	// if additional permissions are passed in, it'll add onto it as well
	createPermissions(tx *gorm.DB, newFile *File, additionalPermissions ...Permission) error

	// Upsert functions updates the users or groups with the new permissions given.
	// will return error if permissions get narrower than parent folder
	upsertUsersPermission(tx *gorm.DB, file *File, permissionNeeded PermissionNeeded, users ...User) error
	upsertGroupsPermission(tx *gorm.DB, file *File, permissionNeeded PermissionNeeded, groups ...Group) error

	// Removes user or groups permissions. Returns error if it becomes narrower
	removeUsersPermission(tx *gorm.DB, file *File, users ...User) error
	removeGroupsPermission(tx *gorm.DB, file *File, users ...Group) error

	// updateMetadataPermission: Ran after upsert, modify on permissions on Files to ensure that the
	// CreatedAt, UpdatedAt, DeletedAt times match with File
	updateMetadataPermission(tx *gorm.DB, file *File) error

	// deleteFilePermissions deletes permission entries when file/folder gets deleted
	deleteFilePermissions(tx *gorm.DB, file *File) error

	// called from User or Group to ensure that user has permissions to file.
	checkGroupHasPermission(tx *gorm.DB, file *File, permissionsNeeded PermissionNeeded, group *Group) error
	checkUserHasPermission(tx *gorm.DB, file *File, permissionsNeeded PermissionNeeded, user *User) error
}
