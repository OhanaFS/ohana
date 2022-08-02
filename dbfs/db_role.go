package dbfs

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Role struct {
	RoleID      int `gorm:"primaryKey;autoIncrement"`
	RoleMapping string
	RoleName    string   `gorm:"not null"`
	Users       []*User  `gorm:"many2many:user_roles;"`
	Groups      []*Group `gorm:"many2many:group_roles;"`
}

// CreateNewRole creates a new role
func CreateNewRole(tx *gorm.DB, name, mapping string) (*Role, error) {
	var role Role
	role.RoleName = name
	role.RoleMapping = mapping
	if err := tx.Create(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// GetRoleByID returns a role by its ID
func GetRoleByID(tx *gorm.DB, id int) (*Role, error) {
	var role Role
	if err := tx.Where("role_id = ?", id).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// GetRoleByName returns a role by its name
func GetRoleByName(tx *gorm.DB, name string) (*Role, error) {
	var role Role
	if err := tx.Preload(clause.Associations).
		Where("role_name = ?", name).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// GetRolesByNames returns a list of roles matching a list of names
func GetRolesByNames(tx *gorm.DB, names []string) ([]Role, error) {
	var roles []Role
	if err := tx.Preload(clause.Associations).
		Where("role_name IN ?", names).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}
