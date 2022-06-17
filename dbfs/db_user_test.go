package dbfs_test

import (
	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/util/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUsers(t *testing.T) {

	db := testutil.NewMockDB(t)

	user1, err := dbfs.CreateNewUser(db, "test1",
		"Test1", dbfs.AccountTypeAdmin, "BLAH1")
	assert.Nil(t, err)

	user2, err := dbfs.CreateNewUser(db, "test2",
		"Test2", dbfs.AccountTypeEndUser, "BLAH2")
	assert.Nil(t, err)

	user3, err := dbfs.CreateNewUser(db, "test3",
		"Test3", dbfs.AccountTypeAdmin, "BLAH3")
	assert.Nil(t, err)

	t.Run("Verifying newly created users", func(t *testing.T) {
		assert := assert.New(t)

		assert.Equal("test1", user1.Username)
		assert.Equal("Test1", user1.Name)
		assert.Equal(dbfs.AccountTypeAdmin, user1.AccountType)
		assert.Equal("BLAH1", user1.MappedID)

		assert.Equal("test2", user2.Username)
		assert.Equal("Test2", user2.Name)
		assert.Equal(dbfs.AccountTypeEndUser, user2.AccountType)
		assert.Equal("BLAH2", user2.MappedID)

		assert.Equal("test3", user3.Username)
		assert.Equal("Test3", user3.Name)
		assert.Equal(dbfs.AccountTypeAdmin, user3.AccountType)
		assert.Equal("BLAH3", user3.MappedID)
	})

	t.Run("Duplicate Username not allowed", func(t *testing.T) {
		assert := assert.New(t)

		_, err := dbfs.CreateNewUser(db, "test1",
			"Test1", dbfs.AccountTypeAdmin, "BLAH1")

		assert.Equal(dbfs.ErrUsernameExists, err)

	})

	t.Run("Get User based on username", func(t *testing.T) {
		assert := assert.New(t)

		user3, err := dbfs.GetUser(db, "test3")

		assert.Nil(err)

		assert.Equal("test3", user3.Username)
	})

	t.Run("Getting invalid username", func(t *testing.T) {
		assert := assert.New(t)

		_, err := dbfs.GetUser(db, "blahblah")

		assert.Equal(dbfs.ErrUserNotFound, err)

	})

	t.Run("Attempt to De-acctivate and activate a user", func(t *testing.T) {
		assert := assert.New(t)

		assert.Equal(true, user1.Activated)

		user1.DeactivateUser(db)

		assert.Equal(false, user1.Activated)

		user1.ActivateUser(db)

		assert.Equal(true, user1.Activated)

	})

}
