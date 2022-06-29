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

	// Create SuperUser if it doesn't already exist
	var superUser *User
	err = db.Where("username = ?", "superuser").First(&superUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			superUser, err = CreateNewUser(db, "superuser", "Super User", 2, "")
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	// Set root folder if non existent
	var rootFolder *File
	err = db.Where("file_id = ?", "00000000-0000-0000-0000-000000000000").First(&rootFolder).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			rootFolder = &File{
				FileId:        "00000000-0000-0000-0000-000000000000",
				FileName:      "root",
				EntryType:     ISFOLDER,
				VersionNo:     0,
				Size:          0,
				ActualSize:    0,
				CreatedTime:   time.Time{},
				ModifiedTime:  time.Time{},
				Status:        1,
				HandledServer: "",
			}
			if err = db.Save(rootFolder).Error; err != nil {
				return err
			}
			permission := Permission{
				FileId:     rootFolder.FileId,
				User:       *superUser,
				UserId:     &superUser.UserId,
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
	}

	return nil

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
