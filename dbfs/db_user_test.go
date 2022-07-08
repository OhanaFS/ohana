package dbfs_test

import (
	"testing"

	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/util/testutil"
	"github.com/stretchr/testify/assert"
)

func TestUsers(t *testing.T) {

	db := testutil.NewMockDB(t)

	user1, err := dbfs.CreateNewUser(db, "test1",
		"Test1", dbfs.AccountTypeAdmin, "BLAH1", "refreshToken", "accessToken", "idToken")
	assert.Nil(t, err)

	user2, err := dbfs.CreateNewUser(db, "test2",
		"Test2", dbfs.AccountTypeEndUser, "BLAH2", "refreshToken", "accessToken", "idToken")
	assert.Nil(t, err)

	user3, err := dbfs.CreateNewUser(db, "test3",
		"Test3", dbfs.AccountTypeAdmin, "BLAH3", "refreshToken", "accessToken", "idToken")
	assert.Nil(t, err)

	t.Run("Verifying newly created users", func(t *testing.T) {
		assert := assert.New(t)

		assert.Equal("test1", user1.Email)
		assert.Equal("Test1", user1.Name)
		assert.Equal(dbfs.AccountTypeAdmin, user1.AccountType)
		assert.Equal("BLAH1", user1.MappedId)

		assert.Equal("test2", user2.Email)
		assert.Equal("Test2", user2.Name)
		assert.Equal(dbfs.AccountTypeEndUser, user2.AccountType)
		assert.Equal("BLAH2", user2.MappedId)

		assert.Equal("test3", user3.Email)
		assert.Equal("Test3", user3.Name)
		assert.Equal(dbfs.AccountTypeAdmin, user3.AccountType)
		assert.Equal("BLAH3", user3.MappedId)
	})

	t.Run("Duplicate Email not allowed", func(t *testing.T) {
		assert := assert.New(t)

		_, err := dbfs.CreateNewUser(db, "test1",
			"Test1", dbfs.AccountTypeAdmin, "BLAH1", "refreshToken", "accessToken", "idToken")

		assert.Equal(dbfs.ErrUsernameExists, err)

	})

	t.Run("Get User based on username", func(t *testing.T) {
		assert := assert.New(t)

		user3, err := dbfs.GetUser(db, "test3")

		assert.Nil(err)

		assert.Equal("test3", user3.Email)
	})

	t.Run("Get User based on userId", func(t *testing.T) {
		assert := assert.New(t)

		user3, err := dbfs.GetUserById(db, user3.UserId)

		assert.Nil(err)

		assert.Equal("test3", user3.Email)
	})

	t.Run("Getting invalid username", func(t *testing.T) {
		assert := assert.New(t)

		_, err := dbfs.GetUser(db, "blahblah")

		assert.Equal(dbfs.ErrUserNotFound, err)

	})

	t.Run("Attempt to De-activate and activate a user", func(t *testing.T) {
		assert := assert.New(t)

		assert.Equal(true, user1.Activated)

		err = user1.DeactivateUser(db)
		assert.Nil(err)
		assert.Equal(false, user1.Activated)

		err = user1.ActivateUser(db)
		assert.Nil(err)
		assert.Equal(true, user1.Activated)

	})

}
