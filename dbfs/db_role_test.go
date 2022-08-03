package dbfs_test

import (
	"testing"

	"golang.org/x/exp/slices"

	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/util/slice"
	"github.com/OhanaFS/ohana/util/testutil"
	"github.com/stretchr/testify/assert"
)

func TestRoles(t *testing.T) {
	db := testutil.NewMockDB(t)

	// Create dummy data
	group1, err := dbfs.CreateNewGroup(db, "Water", "melon")
	assert.Nil(t, err)
	group2, err := dbfs.CreateNewGroup(db, "Bana", "na")
	assert.Nil(t, err)

	role1, err := dbfs.CreateNewRole(db, "Neko", "nyan")
	assert.Nil(t, err)
	role2, err := dbfs.CreateNewRole(db, "Cute", "kawaii")
	assert.Nil(t, err)

	group1.AddMappedRole(db, role1)
	group1.AddMappedRole(db, role2)
	group2.AddMappedRole(db, role2)

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

	// Get multiple roles by name
	t.Run("GetRolesByNames", func(t *testing.T) {
		roles, err := dbfs.GetRolesByNames(db, []string{"Neko", "Cute"})
		assert.Nil(t, err)
		assert.EqualValues(t, 2, len(roles))
		roleNames := slice.Map(roles, func(role dbfs.Role) string {
			return role.RoleName
		})
		assert.True(t, slices.Contains(roleNames, role1.RoleName))
		assert.True(t, slices.Contains(roleNames, role2.RoleName))
	})

	// Get all groups by role name
	t.Run("GetGroupsByRoleNames", func(t *testing.T) {
		groups, err := dbfs.GetGroupsByRoleNames(db, []string{"Neko", "Cute"})
		assert.Nil(t, err)
		assert.EqualValues(t, 2, len(groups))

		groups, err = dbfs.GetGroupsByRoleNames(db, []string{"Neko"})
		assert.Nil(t, err)
		assert.EqualValues(t, 1, len(groups))
	})
}
