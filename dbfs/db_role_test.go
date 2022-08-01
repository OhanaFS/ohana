package dbfs_test

import (
	"testing"

	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/util/testutil"
	"github.com/stretchr/testify/assert"
)

func TestRoles(t *testing.T) {
	db := testutil.NewMockDB(t)

	// Create dummy data
	// group1, err := dbfs.CreateNewGroup(db, "Water", "melon")
	// assert.Nil(t, err)
	// group2, err := dbfs.CreateNewGroup(db, "Bana", "na")
	// assert.Nil(t, err)

	role1, err := dbfs.CreateNewRole(db, "Neko", "nyan")
	assert.Nil(t, err)
	// role2, err := dbfs.CreateNewRole(db, "Cute", "kawaii")
	// assert.Nil(t, err)

	// Get the roles
	t.Run("GetRoleByID", func(t *testing.T) {
		role, err := dbfs.GetRoleByID(db, role1.RoleID)
		assert.Nil(t, err)
		assert.Equal(t, role1.RoleID, role.RoleID)
		assert.Equal(t, role1.RoleName, role.RoleName)
		assert.Equal(t, role1.RoleMapping, role.RoleMapping)
		assert.Equal(t, 0, len(role.Users))
		assert.Equal(t, 0, len(role.Groups))
	})
}
