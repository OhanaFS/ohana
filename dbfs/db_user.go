package dbfs

import (
	"errors"
	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

const (
	AccountTypeEndUser = int8(1)
	AccountTypeAdmin   = int8(2)
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUsernameExists  = errors.New("username already is in use")
	ErrInvalidUserType = errors.New("invalid account type")
)

type User struct {
	UserID      string `gorm:"primaryKey; not null"` // Random UUID
	Name        string
	Username    string `gorm:"not null; unique"`
	MappedID    string `gorm:"not null; unique"`
	LastLogin   time.Time
	Activated   bool           `gorm:"not null; default: true"`
	AccountType int8           `gorm:"not null; default: 1"` // 1 = End User, 2 = Admin
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Groups      []*Group       `gorm:"many2many:user_groups;"`
}

//	Groups      []Group        `gorm:"many2many:user_groups"`

type UserInterface interface {
	ModifyName(tx *gorm.DB, newName string) error
	ModifyUsername(tx *gorm.DB, newUsername string) error
	ModifyAccountType(tx *gorm.DB, newStatus int8) error
	MapToNewAccount(tx *gorm.DB, newID string) error
	GetGroupsWithUser(tx *gorm.DB) ([]Group, error)
	DeactivateUser(tx *gorm.DB) error
	DeleteUser(tx *gorm.DB) error
	ActivateUser(tx *gorm.DB) error
	HasPermission(tx *gorm.DB, file *File, needed *PermissionNeeded) (bool, error)
	AddToGroup(tx *gorm.DB, group *Group) error
}

// Compile time assertion to ensure that User follows UserInterface interface.
var _ UserInterface = &User{}

// CreateNewUser
// Creates a new user with a DB provided.
// Requires username, name, AccountType, MappedID
func CreateNewUser(tx *gorm.DB, username string, name string, accountType int8, mappedID string) (*User, error) {

	// Check stuff like enums

	if !(accountType >= 1 || accountType <= 2) {
		return nil, ErrInvalidUserType
	}

	// Create the account

	userAccount := &User{
		UserID:      uuid.New().String(),
		Name:        name,
		Username:    username,
		MappedID:    mappedID,
		Activated:   true,
		AccountType: accountType,
	}

	result := tx.Create(&userAccount)

	if result.Error != nil {
		if err, ok := result.Error.(sqlite3.Error); ok && err.Code == sqlite3.ErrConstraint {
			return nil, ErrUsernameExists
		}
		return nil, result.Error
	} else {
		return userAccount, nil
	}

}

// GetUser returns the User struct based on the given username
func GetUser(tx *gorm.DB, username string) (*User, error) {

	user := &User{}

	if err := tx.Preload(clause.Associations).First(&user, "username = ?", username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// ErrorHandling
			return user, ErrUserNotFound
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

	err = tx.Select(clause.Associations).Delete(&user).Error

	if err != nil {
		return err
	} else {
		return nil
	}

}

// ModifyName modifies the name of the given User and saves it instantly
func (user *User) ModifyName(tx *gorm.DB, newName string) error {
	user.Name = newName
	return tx.Save(&user).Error
}

// ModifyUsername modifies the username of the given User and saves it instantly
func (user *User) ModifyUsername(tx *gorm.DB, NewUsername string) error {
	user.Username = NewUsername
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
func (user *User) MapToNewAccount(tx *gorm.DB, NewID string) error {
	user.MappedID = NewID
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

	var groups []Group

	err := tx.Preload(clause.Associations).Model(&user).Association("Groups").Find(&groups)

	return groups, err
}

// AddToGroup appends the user to a given group.
func (user *User) AddToGroup(tx *gorm.DB, group *Group) error {

	// Checking to ensure that there's no duplicate association
	for _, existingGroup := range user.Groups {
		if existingGroup.GroupID == group.GroupID {
			return nil
		}
	}

	err := tx.Model(&user).Association("Groups").Append([]Group{*group})
	tx.Save(&user)
	return err
}
