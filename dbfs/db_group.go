package dbfs

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Group struct {
	GroupID       string `gorm:"primaryKey"`
	GroupName     string `gorm:"not null"`
	Activated     bool   `gorm:"not null"`
	MappedGroupID string
	Users         []*User `gorm:"many2many:user_groups;"`
}

type GroupInterface interface {
	ModifyName(tx *gorm.DB, newGroupName string) error
	Activate(tx *gorm.DB) error
	Deactivate(tx *gorm.DB) error
	GetUsers(tx *gorm.DB) ([]User, error)
	ModifyMappedGroupID(tx *gorm.DB, newMappedGroup string) error
	HasPermission(tx *gorm.DB, file *File, needed PermissionNeeded) (bool, error)
}

// CreateNewGroup creates a new group
func CreateNewGroup(tx *gorm.DB, name string, mappedGroupID string) (*Group, error) {

	NewGroup := &Group{
		GroupID:       uuid.New().String(),
		GroupName:     name,
		Activated:     true,
		MappedGroupID: mappedGroupID,
	}

	result := tx.Create(&NewGroup)

	// AFAIK the only possible non system error would be is an identical UUID.
	if result.Error != nil {
		return nil, result.Error
	} else {
		return NewGroup, nil
	}
}

// GetGroupsLikeName returns groups based on "search"
// Does not automatically return users associated.
// To see users associated with Group, use GetUsers()
func GetGroupsLikeName(tx *gorm.DB, groupName string) ([]Group, error) {

	var groups []Group

	err := tx.Where("group_name like ?", "%"+groupName+"%").Find(&groups).Error

	if err != nil {
		return nil, err
	}

	return groups, nil
}

// GetGroupBasedOnGroupID returns the group that matches the ID
// Does not automatically return users associated.
// To see users associated with Group, use GetUsers()
func GetGroupBasedOnGroupID(tx *gorm.DB, groupID string) (*Group, error) {
	var group = &Group{GroupID: groupID}

	err := tx.First(group).Error

	if err != nil {
		return nil, err
	}

	return group, nil
}

// DeleteGroup deletes a group and ensures that all associations between User and Group are deleted.
func DeleteGroup(tx *gorm.DB, group *Group) error {

	// Deletes user links
	err := tx.Select(clause.Associations).Delete(group).Error

	if err != nil {
		return err
	}

	err = tx.Select("Group").Delete(group).Error

	return err
}

// Compile time assertion to ensure that User follows UserInterface interface.
var _ GroupInterface = &Group{}

// ModifyName modifies the name of the given Group and saves it instantly
func (g *Group) ModifyName(tx *gorm.DB, newGroupName string) error {
	g.GroupName = newGroupName
	return tx.Save(&g).Error
}

// Activate instantly activates a group
func (g *Group) Activate(tx *gorm.DB) error {
	g.Activated = true
	return tx.Save(&g).Error
}

// Deactivate instantly activates a group
func (g *Group) Deactivate(tx *gorm.DB) error {
	g.Activated = false
	return tx.Save(&g).Error
}

// GetUsers returns the users associated with a group.
func (g *Group) GetUsers(tx *gorm.DB) ([]User, error) {

	var users []User
	err := tx.Preload(clause.Associations).Model(&g).Association("Users").Find(&users)

	return users, err
}

// ModifyMappedGroupID modifies the mapped group ID and saves it instantly.
func (g *Group) ModifyMappedGroupID(tx *gorm.DB, newMappedGroup string) error {
	g.MappedGroupID = newMappedGroup

	return tx.Save(&g).Error
}
