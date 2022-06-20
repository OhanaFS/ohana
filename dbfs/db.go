package dbfs

import (
	"gorm.io/gorm"
	"time"
)

// InitDB Initiates the DB with gorm.db.AutoMigrate
func InitDB(db *gorm.DB) error {
	err := db.AutoMigrate(&User{}, &Group{}, &File{}, &FileVersion{}, &Fragment{}, &Permission{}, &PermissionHistory{})

	if err != nil {
		return err
	}

	// Create SuperUser

	superUser, err := CreateNewUser(db, "superuser", "Super User", 2, "")

	if err != nil {
		return err
	}

	// Set root folder

	rootFolder := File{
		FileID:        "00000000-0000-0000-0000-000000000000",
		FileName:      "root",
		EntryType:     0,
		VersionNo:     0,
		Size:          0,
		ActualSize:    0,
		CreatedTime:   time.Time{},
		ModifiedTime:  time.Time{},
		Status:        1,
		HandledServer: "",
	}

	if err = db.Save(&rootFolder).Error; err != nil {
		return err
	}

	// Assign superuser permission to root folder

	permission := Permission{
		FileID:     rootFolder.FileID,
		User:       *superUser,
		UserID:     superUser.UserID,
		CanRead:    true,
		CanWrite:   true,
		CanExecute: true,
		CanShare:   true,
		VersionNo:  0,
		Audit:      false,
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
		Status:     1,
	}

	return db.Save(&permission).Error
}

type PermissionNeeded struct {
	Read    bool
	Write   bool
	Execute bool
	Share   bool
	Audit   bool
}

func (current PermissionNeeded) UpdatePermissionsHave(newPermissions Permission) {
	if newPermissions.CanRead {
		current.Read = true
	}
	if newPermissions.CanWrite {
		current.Write = true
	}
	if newPermissions.CanExecute {
		current.Execute = true
	}
	if newPermissions.CanShare {
		current.Share = true
	}
	if newPermissions.Audit {
		current.Audit = true
	}
}

// UpdatePermissionsHaveForSharing either updates current record with superseding permissions
// or returns false
func (current PermissionNeeded) UpdatePermissionsHaveForSharing(newPermissions Permission) bool {
	// Check if any incoming permissions are lower than current permissions
	if (current.Read && !newPermissions.CanRead) || (current.Write && !newPermissions.CanWrite) ||
		(current.Execute && !newPermissions.CanExecute) || (current.Share && !newPermissions.CanShare) ||
		(current.Audit && !newPermissions.Audit) {
		return false
	} else {
		current.UpdatePermissionsHave(newPermissions)
		return true
	}
}

func (current PermissionNeeded) HasPermissions(incomingPermissions PermissionNeeded) bool {

	return !((incomingPermissions.Read && !current.Read) || (incomingPermissions.Write && !current.Write) ||
		(incomingPermissions.Execute && !current.Execute) || (incomingPermissions.Share && !current.Share))
}
