package dbfs

import (
	"errors"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	AccountTypeEndUser int8 = 1
	AccountTypeAdmin   int8 = 2
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUsernameExists  = errors.New("username already is in use")
	ErrInvalidUserType = errors.New("invalid account type")
	ErrCredentials     = errors.New("invalid credentials")
)

type User struct {
	UserId       string         `gorm:"primaryKey; not null; foreignKey:user_id" json:"user_id"` // Random UUId
	Name         string         `json:"name"`
	Email        string         `gorm:"not null; unique" json:"email"` // Maps to email?
	MappedId     string         `gorm:"not null; unique" json:"-"`     // Maps to userID
	RefreshToken string         `json:"-"`
	AccessToken  string         `json:"-"`
	LastLogin    time.Time      `json:"-"`
	Activated    bool           `gorm:"not null; default: true" json:"-"`
	AccountType  int8           `gorm:"not null; default: 1" json:"-"` // 1 = End User, 2 = Admin
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"-"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"-"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Groups       []*Group       `gorm:"many2many:user_groups;" json:"-"`
	HomeFolderId string         `json:"home_folder_id"`
	//Roles        []*Role        `gorm:"many2many:user_roles;" json:"-"`
}

type SharedWithUser struct {
	UserId      string    `gorm:"primaryKey; not null; foreignKey:user_id" json:"user_id"`
	FileId      string    `gorm:"primaryKey; not null; foreignKey:file_id" json:"file_id"`
	DateCreated time.Time `json:"date_created"`
	File        File      `gorm:"foreignKey:file_id"`
}

type SharedWithGroup struct {
	GroupId     string    `gorm:"primaryKey; not null; foreignKey:group_id" json:"group_id"`
	FileId      string    `gorm:"primaryKey; not null; foreignKey:file_id" json:"file_id"`
	DateCreated time.Time `json:"date_created"`
	File        File      `gorm:"foreignKey:file_id" `
}

type FavoriteFileItems struct {
	UserId string `gorm:"primaryKey; not null; foreignKey:user_id" json:"user_id"`
	FileId string `gorm:"primaryKey; not null; foreignKey:file_id" json:"file_id"`
	File   File   `gorm:"foreignKey:file_id" `
}

//	Groups      []Group        `gorm:"many2many:user_groups"`

type UserInterface interface {
	ModifyName(tx *gorm.DB, newName string) error
	ModifyEmail(tx *gorm.DB, newUsername string) error
	ModifyAccountType(tx *gorm.DB, newStatus int8) error
	MapToNewAccount(tx *gorm.DB, newId string) error
	GetGroupsWithUser(tx *gorm.DB) ([]Group, error)
	DeactivateUser(tx *gorm.DB) error
	DeleteUser(tx *gorm.DB) error
	ActivateUser(tx *gorm.DB) error
	HasPermission(tx *gorm.DB, file *File, needed *PermissionNeeded) (bool, error)
	AddToGroup(tx *gorm.DB, group *Group) error
	SetGroups(tx *gorm.DB, groups []Group) error
	GetFavoriteFiles(tx *gorm.DB, start uint) ([]File, error)
	GetSharedWithUser(tx *gorm.DB) ([]File, error)
}

// Compile time assertion to ensure that User follows UserInterface interface.
var _ UserInterface = &User{}

// CreateNewUser creates a new user with a DB provided.
// Requires username, name, AccountType, MappedId
func CreateNewUser(tx *gorm.DB, email string, name string, accountType int8,
	mappedId, refreshToken, accessToken, idToken, server string) (*User, error) {

	// Validate account type
	if !(accountType >= 1 || accountType <= 2) {
		return nil, ErrInvalidUserType
	}

	// Create the account
	userAccount := &User{
		UserId:       uuid.New().String(),
		Name:         name,
		Email:        email,
		MappedId:     mappedId,
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
		Activated:    true,
		AccountType:  accountType,
		HomeFolderId: uuid.New().String(),
	}
	result := tx.Create(&userAccount)

	if result.Error != nil {
		if err, ok := result.Error.(sqlite3.Error); ok &&
			err.Code == sqlite3.ErrConstraint {
			return nil, ErrUsernameExists
		}
		return nil, result.Error
	}

	// Create folder for new users
	err := tx.Transaction(func(tx *gorm.DB) error {

		var UserFolderIDVar string
		UserFolderIDVar = UsersFolderId

		newFolder := &File{
			FileId:             userAccount.HomeFolderId,
			FileName:           userAccount.UserId,
			MIMEType:           "",
			EntryType:          IsFolder,
			ParentFolderFileId: &UserFolderIDVar,
			VersionNo:          0,
			DataId:             "",
			DataIdVersion:      0,
			Size:               0,
			ActualSize:         0,
			CreatedTime:        time.Time{},
			ModifiedTime:       time.Time{},
			VersioningMode:     VersioningOff,
			Status:             1,
			HandledServer:      server,
		}

		if err := tx.Create(&newFolder).Error; err != nil {
			return err
		}

		newPermission := Permission{
			FileId:     newFolder.FileId,
			User:       *userAccount,
			UserId:     &userAccount.UserId,
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

		return tx.Save(&newPermission).Error
	},
	)
	if err != nil {
		// Delete the userAccount and return the error

		tx.Delete(&userAccount)
		return nil, err
	}

	return userAccount, nil
}

// GetUser returns the User struct based on the given username
func GetUser(tx *gorm.DB, username string) (*User, error) {
	user := &User{}

	if err := tx.Preload(clause.Associations).
		First(&user, "email = ?", username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
	}

	return user, nil
}

// GetUsers returns all the users in the DB
func GetUsers(tx *gorm.DB) ([]User, error) {
	users := []User{}

	if err := tx.Preload(clause.Associations).
		Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// GetUserById returns the User struct based on the given userId
func GetUserById(tx *gorm.DB, userId string) (*User, error) {
	user := &User{}

	if err := tx.Preload(clause.Associations).
		First(&user, "user_id = ?", userId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
	}

	return user, nil
}

// GetUserByMappedId returns the User struct based on the given mappedId
func GetUserByMappedId(tx *gorm.DB, mappedId string) (*User, error) {
	user := &User{}

	if err := tx.Preload(clause.Associations).
		First(&user, "mapped_id = ?", mappedId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
	}

	return user, nil
}

// DeleteUser deletes the user and removes the associated links instantly.
func DeleteUser(tx *gorm.DB, username string) error {
	user, err := GetUser(tx, username)

	if err != nil {
		return err
	}

	return tx.Select(clause.Associations).Delete(&user).Error
}

// ModifyName modifies the name of the given User and saves it instantly
func (user *User) ModifyName(tx *gorm.DB, newName string) error {
	user.Name = newName
	return tx.Save(&user).Error
}

// ModifyEmail modifies the username of the given User and saves it instantly
func (user *User) ModifyEmail(tx *gorm.DB, newEmail string) error {
	user.Email = newEmail
	return tx.Save(&user).Error
}

// ModifyAccountType modifies the type of the given User and saves it instantly
func (user *User) ModifyAccountType(tx *gorm.DB, NewStatus int8) error {
	//	Check if Account Type is valid
	if NewStatus >= 1 || NewStatus <= 2 {
		user.AccountType = NewStatus
	} else {
		return ErrInvalidUserType
	}

	return tx.Save(&user).Error
}

// MapToNewAccount modifies the mapping identity of the User and saves it instantly
func (user *User) MapToNewAccount(tx *gorm.DB, NewId string) error {
	user.MappedId = NewId
	return tx.Save(&user).Error
}

// DeactivateUser deactivates the User and saves it instantly
func (user *User) DeactivateUser(tx *gorm.DB) error {
	user.Activated = false
	return tx.Save(&user).Error
}

// DeleteUser deletes the User
func (user *User) DeleteUser(tx *gorm.DB) error {
	if err := tx.Delete(&user).Error; err != nil {
		return err
	}
	return nil
}

// ActivateUser activates the User and saves it instantly
func (user *User) ActivateUser(tx *gorm.DB) error {
	user.Activated = true
	return tx.Save(&user).Error
}

// GetGroupsWithUser refreshes user object with group data.
func (user *User) GetGroupsWithUser(tx *gorm.DB) ([]Group, error) {

	// NEW CODE BELOW

	// Get the roles associated with the user

	//var roles []Role
	var groups []Group
	//err := tx.
	//	Preload("Roles.Groups").
	//	Preload(clause.Associations).
	//	Model(&user).
	//	Association("Roles").
	//	Find(&roles)
	//
	//// Get groups associated with each role.
	//
	//for _, role := range roles {
	//	for _, group := range role.Groups {
	//		groups = append(groups, *group)
	//	}
	//}

	err := tx.Preload(clause.Associations).Model(&user).Association("Groups").Find(&groups)

	return groups, err
}

// AddToGroup appends the user to a given group.
func (user *User) AddToGroup(tx *gorm.DB, group *Group) error {

	// Checking to ensure that there's no duplicate association
	for _, existingGroup := range user.Groups {
		if existingGroup.GroupId == group.GroupId {
			return nil
		}
	}

	err := tx.Model(&user).Association("Groups").Append([]Group{*group})
	tx.Save(&user)
	return err
}

// RefreshGroups will check if there have been any changes in the groups that the user belongs in
// It will ping the server using the refresh token to get a new access token and id token.
func (user *User) RefreshGroups(tx *gorm.DB) error {

	panic("Not implemented")
}

// SetGroups replaces the user's groups with a new list of groups
func (user *User) SetGroups(tx *gorm.DB, groups []Group) error {
	if err := tx.Model(&user).Association("Groups").Replace(groups); err != nil {
		return err
	}
	return tx.Save(&user).Error
}

// GetFavoriteFiles Gets a list of the user's favorite file
func (user *User) GetFavoriteFiles(tx *gorm.DB, start uint) ([]File, error) {

	var FavoriteFileItems []FavoriteFileItems
	if err := tx.Preload(clause.Associations).
		Find(&FavoriteFileItems, "user_id = ?", user.UserId).
		Offset(int(start)).Limit(50).Error; err != nil {
		return nil, err
	}
	files := make([]File, len(FavoriteFileItems))
	for i, item := range FavoriteFileItems {
		files[i] = item.File
	}
	return files, nil
}

// GetFavoriteFileByFileId Returns the favorite file with the given fileId if it exists
func (user *User) GetFavoriteFileByFileId(tx *gorm.DB, fileId string) (*File, error) {

	var favItem FavoriteFileItems
	if err := tx.Preload(clause.Associations).
		Where("user_id = ? AND file_id = ?", user.UserId, fileId).
		First(&favItem).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrFileNotFound
		}
		return nil, err
	}
	return &favItem.File, nil

}

func (user *User) GetSharedWithUser(tx *gorm.DB) ([]File, error) {

	// First, we will grab all the files that are shared with the user
	// We will store it in a hashmap, with fileId as the key and the date as value
	// This way, we can return the files in order of the date shared
	var SharedWithUsers []SharedWithUser

	if err := tx.Where("user_id = ?", user.UserId).Find(&SharedWithUsers).Error; err != nil {
		return nil, err
	}

	hasAccess := make(map[string]*time.Time)
	for _, item := range SharedWithUsers {
		hasAccess[item.FileId] = &item.DateCreated
	}

	// Now, we will grab all the files that the user has access to via groups
	groups, err := user.GetGroupsWithUser(tx)
	if err != nil {
		return nil, err
	}

	groupsArray := make([]string, len(groups))
	for i, group := range groups {
		groupsArray[i] = group.GroupId
	}

	var SharedWithGroups []SharedWithGroup

	if err := tx.Where("group_id IN (?)", groupsArray).Find(&SharedWithGroups).Error; err != nil {
		return nil, err
	}

	for _, item := range SharedWithGroups {
		value, ok := hasAccess[item.FileId]
		if ok {
			// update if newer date
			if item.DateCreated.After(*value) {
				hasAccess[item.FileId] = &item.DateCreated
			}
		} else {
			hasAccess[item.FileId] = &item.DateCreated
		}
	}

	// Now, we'll sort hasAccess by the value (date)
	keys := make([]string, 0, len(hasAccess))

	for key := range hasAccess {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return hasAccess[keys[i]].Unix() < hasAccess[keys[j]].Unix()
	})

	// keys is now sorted

	files := make([]File, len(keys))
	for i, key := range keys {
		var tempFile File
		if err := tx.Where("file_id = ?", key).First(&tempFile).Error; err != nil {
			return nil, err
		}
		files[i] = tempFile
	}

	return files, nil

}
