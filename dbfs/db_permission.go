package dbfs

import (
	"errors"
	"gorm.io/gorm/clause"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrPermissionsAreNarrowerThanParent = errors.New("new permissions lower than parent permissions")
)

type Permission struct {
	FileId       string         `gorm:"primaryKey; not_null" json:"file_id"`
	PermissionId string         `gorm:"primaryKey;" json:"permission_id"`
	User         User           `gorm:"foreignKey:UserId; references:UserId"`
	UserId       *string        `json:"user_id"`
	Group        Group          `gorm:"foreignKey:GroupId; references:GroupId"`
	GroupId      *string        `json:"group_id"`
	CanRead      bool           `json:"can_read"`
	CanWrite     bool           `json:"can_write"`
	CanExecute   bool           `json:"can_execute"`
	CanShare     bool           `json:"can_share"`
	VersionNo    int            `json:"version_no"`
	Audit        bool           `json:"audit"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Status       int8           `json:"status"` // 1 indicates active, 0 indicates time needs to be updated properly
}

type PermissionHistory struct {
	FileId       string `gorm:"primaryKey"`
	VersionNo    int    `gorm:"primaryKey"`
	PermissionId string `gorm:"primaryKey"`
	User         User   `gorm:"foreignKey:UserId"`
	UserId       *string
	Group        Group `gorm:"foreignKey:GroupId"`
	GroupId      *string
	CanRead      bool
	CanWrite     bool
	CanExecute   bool
	CanShare     bool
	Audit        bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Status       int8
}

// CreatePermissions takes in the parent file and copies it to the new File
// if additional permissions are passed in, it'll add onto it as well
func CreatePermissions(tx *gorm.DB, newFile *File, additionalPermissions ...Permission) error {

	var oldPermissionRecords []Permission

	tx.Where("file_id = ?", newFile.ParentFolderFileId).Find(&oldPermissionRecords)

	// Based on additional permissions, we'll either modify it, or add a new entry

	for _, newPermission := range additionalPermissions {
		newIsUser := newPermission.UserId != nil
		for _, existingPermissions := range oldPermissionRecords {
			existingIsUser := existingPermissions.UserId != nil
			updated := false
			if newIsUser == existingIsUser { // matching type
				if newIsUser { // user
					if newPermission.UserId == existingPermissions.UserId {
						updated = true
					}
				} else { // group
					if newPermission.GroupId == existingPermissions.GroupId {
						updated = true
					}
				}
			}
			if updated {
				// Modify existing permission
				// ONLY ALLOW MORE, NOT LESS
				if !existingPermissions.CanRead {
					existingPermissions.CanRead = newPermission.CanRead
				}

				if !existingPermissions.CanWrite {
					existingPermissions.CanWrite = newPermission.CanWrite
				}

				if !existingPermissions.CanExecute {
					existingPermissions.CanExecute = newPermission.CanExecute
				}

				if !existingPermissions.CanShare {
					existingPermissions.CanShare = newPermission.CanShare
				}

				if !existingPermissions.Audit {
					existingPermissions.Audit = newPermission.Audit
				}
			} else {
				// Add new permission
				oldPermissionRecords = append(oldPermissionRecords, newPermission)
			}
		}
	}

	// Copy all permissions to new file

	for _, permission := range oldPermissionRecords {

		newRecord := Permission{
			FileId:       newFile.FileId,
			PermissionId: uuid.New().String(),
			UserId:       permission.UserId,
			GroupId:      permission.GroupId,
			CanRead:      permission.CanRead,
			CanWrite:     permission.CanWrite,
			CanExecute:   permission.CanExecute,
			CanShare:     permission.CanShare,
			Audit:        permission.Audit,
			VersionNo:    newFile.VersionNo,
			CreatedAt:    newFile.CreatedTime,
			UpdatedAt:    newFile.CreatedTime,
			Status:       1,
		}

		if err := tx.Create(&newRecord).Error; err != nil {
			return err
		}

	}

	return nil
}

// upsertUsersPermission functions updates the users with the new permissions given.
// will return error if permissions get narrower than parent folder
func upsertUsersPermission(tx *gorm.DB, file *File, permissionNeeded *PermissionNeeded, requestUser *User, users ...User) error {

	isFileOrEmptyFolder, err := file.IsFileOrEmptyFolder(tx, requestUser)
	if err != nil {
		return err
	}

	err = tx.Transaction(func(tx *gorm.DB) error {
		if !isFileOrEmptyFolder {
			// recursively go down

			ls, err2 := file.ListContents(tx, requestUser)

			if err2 != nil {
				return err2
			}

			for _, lsFile := range ls {
				if err := upsertUsersPermission(tx, &lsFile, permissionNeeded, requestUser, users...); err != nil {
					return err
				}
			}
		}

		var oldPermissionRecords []Permission

		tx.Where("file_id = ?", file.FileId).Find(&oldPermissionRecords)

		for _, user := range users {
			updated := false
			for i, existingPermissions := range oldPermissionRecords {
				if existingPermissions.UserId != nil {
					if user.UserId == *existingPermissions.UserId {
						updated = true
					}
				}

				if updated {

					// Getting parent file/folder
					var parentFile File
					tx.Where("file_id = ?", file.ParentFolderFileId).First(&parentFile)

					// Only allows removal of permissions if the parent doesn't have it as well.

					// CanRead
					if !existingPermissions.CanRead && permissionNeeded.Read {
						oldPermissionRecords[i].CanRead = true
					} else if existingPermissions.CanRead && !permissionNeeded.Read {
						// Check if parent has permission
						hasPermission, err := user.HasPermission(tx, &parentFile, &PermissionNeeded{Read: true})

						if err != nil {
							if !errors.Is(err, ErrFileNotFound) {
								return err
							}
						}

						if hasPermission {
							return ErrPermissionsAreNarrowerThanParent
						} else {
							oldPermissionRecords[i].CanRead = false
						}
					}

					// CanWrite
					if !existingPermissions.CanWrite && permissionNeeded.Write {
						oldPermissionRecords[i].CanWrite = true
					} else if existingPermissions.CanWrite && !permissionNeeded.Write {
						// Check if parent has permission
						hasPermission, err := user.HasPermission(tx, &parentFile, &PermissionNeeded{Write: true})

						if err != nil {
							if !errors.Is(err, ErrFileNotFound) {
								return err
							}
						}

						if hasPermission {
							return ErrPermissionsAreNarrowerThanParent
						} else {
							oldPermissionRecords[i].CanWrite = false
						}
					}

					// CanExecute
					if !existingPermissions.CanExecute && permissionNeeded.Execute {
						oldPermissionRecords[i].CanExecute = true
					} else if existingPermissions.CanExecute && !permissionNeeded.Execute {
						// Check if parent has permission
						hasPermission, err := user.HasPermission(tx, &parentFile, &PermissionNeeded{Execute: true})

						if err != nil {
							if !errors.Is(err, ErrFileNotFound) {
								return err
							}
						}

						if hasPermission {
							return ErrPermissionsAreNarrowerThanParent
						} else {
							oldPermissionRecords[i].CanExecute = false
						}
					}

					// CanShare
					if !existingPermissions.CanShare && permissionNeeded.Share {
						oldPermissionRecords[i].CanShare = true
					} else if existingPermissions.CanShare && !permissionNeeded.Share {
						// Check if parent has permission
						hasPermission, err := user.HasPermission(tx, &parentFile, &PermissionNeeded{Share: true})

						if err != nil {
							if !errors.Is(err, ErrFileNotFound) {
								return err
							}
						}

						if hasPermission {
							return ErrPermissionsAreNarrowerThanParent
						} else {
							oldPermissionRecords[i].CanShare = false
						}
					}

					// Audit
					if !existingPermissions.Audit && permissionNeeded.Audit {
						oldPermissionRecords[i].Audit = true
					} else if existingPermissions.Audit && !permissionNeeded.Audit {
						// Check if parent has permission
						hasPermission, err := user.HasPermission(tx, &parentFile, &PermissionNeeded{Audit: true})

						if err != nil {
							if !errors.Is(err, ErrFileNotFound) {
								return err
							}
						}

						if hasPermission {
							return ErrPermissionsAreNarrowerThanParent
						} else {
							oldPermissionRecords[i].Audit = false
						}
					}

					break

				}

			}

			if !updated {
				// Append permission
				newPermission := Permission{
					FileId:       file.FileId,
					PermissionId: uuid.New().String(),
					UserId:       &user.UserId,
					CanRead:      permissionNeeded.Read,
					CanWrite:     permissionNeeded.Write,
					CanExecute:   permissionNeeded.Execute,
					CanShare:     permissionNeeded.Share,
					Audit:        permissionNeeded.Audit,
					VersionNo:    file.VersionNo,
					Status:       1,
				}

				oldPermissionRecords = append(oldPermissionRecords, newPermission)
			}

		}

		// Save all new permissions
		err := tx.Transaction(func(tx2 *gorm.DB) error {
			for _, permission := range oldPermissionRecords {
				permission.VersionNo = file.VersionNo
				permission.UpdatedAt = file.ModifiedTime
				if err := tx2.Save(&permission).Error; err != nil {
					return err
				}
			}

			for _, user := range users {
				err := tx.Clauses(clause.OnConflict{DoNothing: true}).Save(&SharedWithUser{
					FileId:      file.FileId,
					UserId:      user.UserId,
					DateCreated: time.Now(),
				}).Error
				if err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			return err
		}

		return updatePermissionsVersions(tx, file, requestUser)

	})

	return nil

}

// upsertGroupsPermission functions updates the groups with the new permissions given.
// will return error if permissions get narrower than parent folder
func upsertGroupsPermission(tx *gorm.DB, file *File, permissionNeeded *PermissionNeeded, requestUser *User, groups ...Group) error {

	isFileOrEmptyFolder, err := file.IsFileOrEmptyFolder(tx, requestUser)
	if err != nil {
		return err
	}

	err = tx.Transaction(func(tx *gorm.DB) error {

		if !isFileOrEmptyFolder {
			// recursively go down
			ls, err2 := file.ListContents(tx, requestUser)

			if err2 != nil {
				return err2
			}

			for _, lsFile := range ls {
				if err2 := upsertGroupsPermission(tx, &lsFile, permissionNeeded, requestUser, groups...); err2 != nil {
					return err2
				}
			}
		}

		var oldPermissionRecords []Permission

		tx.Where("file_id = ?", file.FileId).Find(&oldPermissionRecords)

		for _, group := range groups {
			updated := false
			for i, existingPermissions := range oldPermissionRecords {
				if existingPermissions.GroupId != nil {
					if group.GroupId == *existingPermissions.GroupId {
						updated = true
						// Modify existing permission
						// ONLY ALLOW MORE, NOT LESS
					}
				}

				if updated {
					// Getting parent file/folder
					var parentFile File
					tx.Where("file_id = ?", file.ParentFolderFileId).First(&parentFile)

					// Only allows removal of permissions if the parent doesn't have it as well.

					// CanRead
					if !existingPermissions.CanRead && permissionNeeded.Read {
						oldPermissionRecords[i].CanRead = true
					} else if existingPermissions.CanRead && !permissionNeeded.Read {
						// Check if parent has permission
						hasPermission, err2 := group.HasPermission(tx, &parentFile, &PermissionNeeded{Read: true})

						if err2 != nil {
							if !errors.Is(err2, ErrFileNotFound) {
								return err
							}
						}

						if hasPermission {
							return ErrPermissionsAreNarrowerThanParent
						} else {
							oldPermissionRecords[i].CanRead = false
						}
					}

					// CanWrite
					if !existingPermissions.CanWrite && permissionNeeded.Write {
						oldPermissionRecords[i].CanWrite = true
					} else if existingPermissions.CanWrite && !permissionNeeded.Write {
						// Check if parent has permission
						hasPermission, err2 := group.HasPermission(tx, &parentFile, &PermissionNeeded{Write: true})

						if err2 != nil {
							if !errors.Is(err2, ErrFileNotFound) {
								return err
							}
						}

						if hasPermission {
							return ErrPermissionsAreNarrowerThanParent
						} else {
							oldPermissionRecords[i].CanWrite = false
						}
					}

					// CanExecute
					if !existingPermissions.CanExecute && permissionNeeded.Execute {
						oldPermissionRecords[i].CanExecute = true
					} else if existingPermissions.CanExecute && !permissionNeeded.Execute {
						// Check if parent has permission
						hasPermission, err2 := group.HasPermission(tx, &parentFile, &PermissionNeeded{Execute: true})

						if err2 != nil {
							if !errors.Is(err2, ErrFileNotFound) {
								return err
							}
						}

						if hasPermission {
							return ErrPermissionsAreNarrowerThanParent
						} else {
							oldPermissionRecords[i].CanExecute = false
						}
					}

					// CanShare
					if !existingPermissions.CanShare && permissionNeeded.Share {
						oldPermissionRecords[i].CanShare = true
					} else if existingPermissions.CanShare && !permissionNeeded.Share {
						// Check if parent has permission
						hasPermission, err2 := group.HasPermission(tx, &parentFile, &PermissionNeeded{Share: true})

						if err2 != nil {
							if !errors.Is(err2, ErrFileNotFound) {
								return err
							}
						}

						if hasPermission {
							return ErrPermissionsAreNarrowerThanParent
						} else {
							oldPermissionRecords[i].CanShare = false
						}
					}

					// Audit
					if !existingPermissions.Audit && permissionNeeded.Audit {
						oldPermissionRecords[i].Audit = true
					} else if existingPermissions.Audit && !permissionNeeded.Audit {
						// Check if parent has permission
						hasPermission, err2 := group.HasPermission(tx, &parentFile, &PermissionNeeded{Audit: true})

						if err2 != nil {
							if !errors.Is(err2, ErrFileNotFound) {
								return err
							}
						}

						if hasPermission {
							return ErrPermissionsAreNarrowerThanParent
						} else {
							oldPermissionRecords[i].Audit = false
						}
					}

					break
				}
			}
			if !updated {
				// Append permission
				newPermission := Permission{
					FileId:       file.FileId,
					PermissionId: uuid.New().String(),
					GroupId:      &group.GroupId,
					CanRead:      permissionNeeded.Read,
					CanWrite:     permissionNeeded.Write,
					CanExecute:   permissionNeeded.Execute,
					CanShare:     permissionNeeded.Share,
					Audit:        permissionNeeded.Audit,
					VersionNo:    file.VersionNo,
					Status:       1,
				}

				oldPermissionRecords = append(oldPermissionRecords, newPermission)

			}
		}

		// Save all new permissions
		err := tx.Transaction(func(tx2 *gorm.DB) error {
			for _, permission := range oldPermissionRecords {
				permission.VersionNo = file.VersionNo
				permission.UpdatedAt = file.ModifiedTime
				if err := tx2.Save(&permission).Error; err != nil {
					return err
				}
			}

			for _, group := range groups {
				if tx2.Clauses(clause.OnConflict{DoNothing: true}).Save(&SharedWithGroup{
					FileId:      file.FileId,
					GroupId:     group.GroupId,
					DateCreated: time.Now(),
				}).Error != nil {
					return nil
				}
			}

			return nil
		})
		if err != nil {
			return err
		}

		return updatePermissionsVersions(tx, file, requestUser)
	})

	return err

}

// deleteFilePermissions deletes permission entries when file/folder gets deleted
func deleteFilePermissions(tx *gorm.DB, file *File) error {
	if err := tx.Where("file_id = ?", file.FileId).Delete(&Permission{}).Error; err != nil {
		return err
	}
	return nil
}

// revokeUsersPermission revokes users permissions. Returns error if it becomes narrower
// DOES NOT TAKE INTO ACCOUNT IF THE USER IS IN A GROUP THAT HAS THE PERMISSION
func revokeUsersPermission(tx *gorm.DB, file *File, users []User, invoker *User) error {

	err := tx.Transaction(func(tx *gorm.DB) error {
		for _, user := range users {

			// Check if user has permission in the parent folder (if they do, can't revoke)
			_, err := GetFileById(tx, *file.ParentFolderFileId, &user)

			if err != nil {
				if errors.Is(err, ErrFileNotFound) {
					// Can revoke safely

					var tempPerm Permission

					err = tx.Where(&Permission{
						FileId: file.FileId,
						UserId: &user.UserId,
					}).First(&tempPerm).Error

					err = tx.Delete(&tempPerm).Error
					if err != nil {
						return err
					} else {
						continue
					}
				} else {
					return err
				}
			}

			// if no error, they have permission to the parent folder. Apply the parent permission to the below file
			parentPermission := Permission{}
			err = tx.Where("file_id = ? AND user_id = ?", file.ParentFolderFileId, user.UserId).First(&parentPermission).Error
			if err != nil {
				return err
			}

			// Update current permission
			err = tx.Model(&Permission{}).Where("file_id = ? AND user_id = ?", file.FileId, user.UserId).
				Updates(Permission{
					CanRead:    parentPermission.CanRead,
					CanWrite:   parentPermission.CanWrite,
					CanExecute: parentPermission.CanExecute,
					CanShare:   parentPermission.CanShare,
					Audit:      parentPermission.Audit}).Error

			if err != nil {
				return err
			}

			// Removing from SharedWithUser
			err = tx.Where("file_id = ? AND user_id = ?",
				file.FileId, user.UserId).
				Delete(&SharedWithUser{}).Error
			if err != nil {
				return err
			}

		}

		return updatePermissionsVersions(tx, file, invoker)
	})

	return err

}

// revokeGroupsPermission revokes a Group's permissions. Returns error if it becomes narrower
// DOES NOT TAKE INTO ACCOUNT IF THE USER IS IN A GROUP THAT HAS THE PERMISSION
func revokeGroupsPermission(tx *gorm.DB, file *File, groups []Group) error {

	err := tx.Transaction(func(tx *gorm.DB) error {
		for _, group := range groups {

			// Check if user has permission in the parent folder (if they do, can't revoke)

			var count int64

			err := tx.Where("file_id = ? AND group_id = ?", file.ParentFolderFileId, group.GroupId).Count(&count).Error
			if err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					return err
				}
				count = 0
			}

			if count > 0 {
				return ErrPermissionsAreNarrowerThanParent
			}

			// Check if group has permission in the file and if so delete it.

			err = tx.Where("file_id = ? AND group_id = ?", file.FileId, group.GroupId).Delete(&Permission{}).Error
			if err != nil {
				return err
			}

			err = tx.Where("file_id = ? AND group_id = ?",
				file.FileId, group.GroupId).
				Delete(&SharedWithGroup{}).Error
			if err != nil {
				return err
			}
		}
		return updatePermissionsVersions(tx, file, nil)

	})

	return err

}

// updatePermissionsVersions updates the File, FileVersion, Permissions Version Number
func updatePermissionsVersions(tx *gorm.DB, file *File, user *User) error {
	return tx.Transaction(func(tx *gorm.DB) error {

		// Update File
		file.VersionNo = file.VersionNo + 1
		err2 := tx.Save(&file).Error
		if err2 != nil {
			return err2
		}

		// Update FileVersion
		err2 = CreateFileVersionFromFile(tx, file, user)
		if err2 != nil {
			return err2
		}

		// Updating current permission records to the new version number
		err2 = tx.Model(&Permission{}).Where("file_id = ?", file.FileId).Updates(map[string]interface{}{"version_no": file.VersionNo}).Error
		if err2 != nil {
			return err2
		}

		// Dumping new permissions to PermissionHistory

		var permissions []Permission
		err2 = tx.Where("file_id = ?", file.FileId).Find(&permissions).Error
		if err2 != nil {
			return err2
		}

		for _, permission := range permissions {
			err2 = tx.Save(&PermissionHistory{
				FileId:       permission.FileId,
				PermissionId: permission.PermissionId,
				VersionNo:    permission.VersionNo,
				GroupId:      permission.GroupId,
				UserId:       permission.UserId,
				CanRead:      permission.CanRead,
				CanWrite:     permission.CanWrite,
				CanExecute:   permission.CanExecute,
				CanShare:     permission.CanShare,
				Audit:        permission.Audit,
				Status:       permission.Status,
				CreatedAt:    permission.CreatedAt,
				UpdatedAt:    permission.UpdatedAt,
			}).Error
			if err2 != nil {
				return err2
			}
		}

		return nil
	})

}

// updateMetadataPermission updates the version No and updated time to the File updated date.
func updateMetadataPermission(tx *gorm.DB, file *File) error {

	return tx.Where("file_id = ?", file.FileId).Updates(map[string]interface{}{"version_no": file.VersionNo,
		"updated_at": file.ModifiedTime}).Error

}

// GetPermissionHistory returns the PermissionHistory for a file
func GetPermissionHistory(tx *gorm.DB, file *File, user *User) ([]PermissionHistory, error) {

	// Check if the person has permission to view the file
	_, err := GetFileById(tx, file.FileId, user)
	if err != nil {
		return nil, err
	}

	var permissions []PermissionHistory
	err = tx.Model(&Permission{}).Where("file_id = ?", file.FileId).Find(&permissions).Error
	return permissions, err
}

// HasPermission verifies that the user has the permission requested to a file.
func (user *User) HasPermission(tx *gorm.DB, file *File, needed *PermissionNeeded) (bool, error) {

	if user.AccountType == AccountTypeAdmin {
		return true, nil
	}

	hasPermission := false

	// Sharing permissions are separate from non-sharing permissions
	noSharePermissions := PermissionNeeded{}
	var sharePermissions []PermissionNeeded

	var permission Permission

	if err := tx.Where("file_id = ? AND user_id = ?", file.FileId, user.UserId).First(&permission).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			groups, err2 := user.GetGroupsWithUser(tx)
			if err2 != nil {
				if errors.Is(err2, gorm.ErrRecordNotFound) {
					return false, ErrFileNotFound
				}
				return false, err
			}

			for _, group := range groups {

				var permission Permission

				if err = tx.Where("file_id = ? AND group_id = ?", file.FileId, group.GroupId).Find(&permission).Error; err != nil {
					return false, err
				}

				noSharePermissions.UpdatePermissionsHave(permission)

				if permission.CanShare && needed.Share {
					if len(sharePermissions) == 0 {
						sharePermissions = append(sharePermissions, PermissionNeeded{
							Read:    permission.CanRead,
							Write:   permission.CanWrite,
							Execute: permission.CanExecute,
							Share:   permission.CanShare,
							Audit:   permission.Audit,
						})
					}

					updated := false

					for _, sharePermission := range sharePermissions {
						updated = sharePermission.UpdatePermissionsHaveForSharing(permission)
						if updated {
							break
						}
					}
					if !updated {
						sharePermissions = append(sharePermissions, PermissionNeeded{
							Read:    permission.CanRead,
							Write:   permission.CanWrite,
							Execute: permission.CanExecute,
							Share:   permission.CanShare,
							Audit:   permission.Audit,
						})
					}
				}

				if needed.Share {
					for _, sharePermission := range sharePermissions {
						if needed.HasPermissions(sharePermission) {
							return true, nil
						}
					}
				} else {
					if needed.HasPermissions(noSharePermissions) {
						return true, nil
					}
				}

			}

		} else {
			return false, err
		}
	} else {
		hasPermission = true

		// See if the user has permissions to it

		if needed.Read && !permission.CanRead {
			hasPermission = false
		}
		if needed.Write && !permission.CanWrite {
			hasPermission = false
		}
		if needed.Execute && !permission.CanExecute {
			hasPermission = false
		}
		if needed.Share && !permission.CanShare {
			hasPermission = false
		}
		if needed.Audit && !permission.Audit {
			hasPermission = false
		}
	}

	return hasPermission, nil

}

// HasPermission will verify that a group has the permission required.
func (g *Group) HasPermission(tx *gorm.DB, file *File, needed *PermissionNeeded) (bool, error) {
	// Get all permissions for the group for that file

	var permission Permission

	if err := tx.Where("file_id = ? AND group_id = ?", file.FileId, g.GroupId).Find(&permission).Error; err != nil {
		return false, err
	}

	// Check if the permissions are sufficient

	hasPermission := true

	if needed.Read && !permission.CanRead {
		hasPermission = false
	}
	if needed.Write && !permission.CanWrite {
		hasPermission = false
	}
	if needed.Execute && !permission.CanExecute {
		hasPermission = false
	}
	if needed.Share && !permission.CanShare {
		hasPermission = false
	}
	if needed.Audit && !permission.Audit {
		hasPermission = false
	}

	return hasPermission, nil

}
