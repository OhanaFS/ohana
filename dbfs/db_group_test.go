package dbfs_test

import (
	"fmt"
	"testing"

	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/util/slice"
	"github.com/OhanaFS/ohana/util/testutil"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

func TestGroups(t *testing.T) {

	db := testutil.NewMockDB(t)

	// Creating Groups to test with
	group1, err := dbfs.CreateNewGroup(db, "PogGroup1", "ANewMapping1")
	assert.Nil(t, err)

	group2, err := dbfs.CreateNewGroup(db, "PogGroup2", "ANewMapping2")
	assert.Nil(t, err)

	group3, err := dbfs.CreateNewGroup(db, "PogGroup3", "ANewMapping3")
	assert.Nil(t, err)

	// Creating Users to test with
	user1, err := dbfs.CreateNewUser(db, "test1",
		"Test1", dbfs.AccountTypeAdmin, "BLAH1", "refreshToken", "accessToken",
		"idToken", "server")
	assert.Nil(t, err)

	user2, err := dbfs.CreateNewUser(db, "test2",
		"Test2", dbfs.AccountTypeEndUser, "BLAH2", "refreshToken", "accessToken",
		"idToken", "server")
	assert.Nil(t, err)

	user3, err := dbfs.CreateNewUser(db, "test3",
		"Test3", dbfs.AccountTypeAdmin, "BLAH3", "refreshToken", "accessToken",
		"idToken", "server")
	assert.Nil(t, err)

	user4, err := dbfs.CreateNewUser(db, "test4",
		"Test4", dbfs.AccountTypeEndUser, "BLAH4", "refreshToken", "accessToken",
		"idToken", "server")
	assert.Nil(t, err)

	_, err = dbfs.CreateNewUser(db, "test5",
		"Test5", dbfs.AccountTypeEndUser, "BLAH5", "refreshToken", "accessToken",
		"idToken", "server")
	assert.Nil(t, err)

	var testGroupId string

	t.Run("Verifying newly created groups", func(t *testing.T) {

		assert := assert.New(t)

		assert.Equal("PogGroup1", group1.GroupName)
		assert.Equal("PogGroup2", group2.GroupName)
		assert.Equal("PogGroup3", group3.GroupName)

		testGroupId = group1.GroupId

		assert.Equal("ANewMapping1", group1.MappedGroupId)
		assert.Equal("ANewMapping2", group2.MappedGroupId)
		assert.Equal("ANewMapping3", group3.MappedGroupId)

	})

	t.Run("Verifying that adding users to groups works", func(t *testing.T) {
		assert := assert.New(t)

		// Adding User 1, 2 to group 1
		err := user1.AddToGroup(db, group1)
		assert.Nil(err)
		err = user2.AddToGroup(db, group1)
		assert.Nil(err)

		// Adding User 2, 3 to group 2
		err = user2.AddToGroup(db, group2)
		assert.Nil(err)
		err = user3.AddToGroup(db, group2)
		assert.Nil(err)

		// Adding User 4 to group 3
		err = user4.AddToGroup(db, group3)
		assert.Nil(err)

		// Verifying that groups contain the users.
		// Group 1
		users, err := group1.GetUsers(db)
		assert.Nil(err)

		UserIds := user1.UserId + user2.UserId
		for _, u := range users {
			assert.Contains(UserIds, u.UserId)
		}

		// Group 2
		users, err = group2.GetUsers(db)
		assert.Nil(err)

		UserIds = user2.UserId + user3.UserId
		for _, u := range users {
			assert.Contains(UserIds, u.UserId)
		}

		// Group 3
		users, err = group3.GetUsers(db)
		assert.Nil(err)

		UserIds = user4.UserId
		for _, u := range users {
			assert.Contains(UserIds, u.UserId)
		}

	})

	t.Run("Verifying group replacement works", func(t *testing.T) {
		assert := assert.New(t)

		test := func(user *dbfs.User, groups []dbfs.Group) {
			db.Transaction(func(tx *gorm.DB) error {
				groupIds := slice.Map(groups, func(g dbfs.Group) string { return g.GroupId })

				// Set user's groups
				err := user1.SetGroups(tx, groups)
				assert.Nil(err)

				// Verify
				newGroups, err := user1.GetGroupsWithUser(tx)
				assert.Nil(err)
				assert.EqualValues(len(groups), len(newGroups), "Group lengths differ")
				for _, newGroup := range newGroups {
					assert.True(slices.Contains(groupIds, newGroup.GroupId), "User has groups outside the replacement set")
				}

				// Cancel the transaction
				return fmt.Errorf("OK")
			})
		}

		test(user1, []dbfs.Group{})
		test(user1, []dbfs.Group{*group3, *group1, *group2})
		test(user1, []dbfs.Group{*group1, *group2})
	})

	// It's presumed at this stage User 1 is already associated with group 1
	t.Run("Verifying that adding a user to a group multiple time still works fine", func(t *testing.T) {

		assert := assert.New(t)

		err = user1.AddToGroup(db, group1)
		assert.Nil(err)

		// Verifying that groups contain the users.
		// Group 1
		users, err := group1.GetUsers(db)
		assert.Nil(err)

		UserIds := user1.UserId + user2.UserId
		for _, u := range users {
			assert.Contains(UserIds, u.UserId)
		}

	})

	t.Run("Verifying that querying a user without a group works", func(t *testing.T) {
		assert := assert.New(t)

		testUser, err := dbfs.GetUser(db, "test5")
		assert.Nil(err)

		groups, err := testUser.GetGroupsWithUser(db)

		assert.Equal(len(groups), 0)

	})

	t.Run("Modifying name", func(t *testing.T) {
		assert := assert.New(t)

		err = group1.ModifyName(db, "PoggiesGroup1")
		assert.Nil(err)

		assert.Equal("PoggiesGroup1", group1.GroupName)

	})

	t.Run("Modifying MappedId", func(t *testing.T) {
		assert := assert.New(t)

		err = group1.ModifyMappedGroupId(db, "NEWMAPPING1")
		assert.Nil(err)

		assert.Equal("NEWMAPPING1", group1.MappedGroupId)

	})

	t.Run("Activation/Deactivation of a Group", func(t *testing.T) {
		assert := assert.New(t)

		err = group1.Deactivate(db)
		assert.Nil(err)

		assert.Equal(false, group1.Activated)

		err = group1.Activate(db)
		assert.Nil(err)

		assert.Equal(true, group1.Activated)

	})

	// Ensuring that deleting a group will update things properly.
	t.Run("Deleting a group", func(t *testing.T) {
		assert := assert.New(t)

		err := dbfs.DeleteGroup(db, group3)

		assert.Nil(err)

		user4, err = dbfs.GetUser(db, "test4")

		assert.Equal(0, len(user4.Groups))

	})

	t.Run("Deleting User updates GetUsers()", func(t *testing.T) {
		assert := assert.New(t)

		err := dbfs.DeleteUser(db, "test2")
		assert.Nil(err)

		// Update group
		group2, err := dbfs.GetGroupBasedOnGroupId(db, group2.GroupId)
		assert.Nil(err)

		users, err := group2.GetUsers(db)
		assert.Nil(err)
		assert.Equal(1, len(users))

	})

	t.Run("Trying to find groups", func(t *testing.T) {

		assert := assert.New(t)

		groups, err := dbfs.GetGroupsLikeName(db, "Group")

		assert.Nil(err)

		assert.Equal(2, len(groups))

	})

	t.Run("Getting a group by by Id", func(t *testing.T) {

		assert := assert.New(t)

		testGroup, err := dbfs.GetGroupBasedOnGroupId(db, testGroupId)

		assert.Nil(err)

		assert.Equal(testGroupId, testGroup.GroupId)

	})
}
