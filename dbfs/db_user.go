package dbfs

import (
	"errors"
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
