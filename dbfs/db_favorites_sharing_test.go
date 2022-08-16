package dbfs_test

import (
	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/util/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestSharingFavorite is a test to see if favorites work and
// if we share a file between users, it should pick up in GetSharedWithUser
func TestSharingFavorite(t *testing.T) {

	db := testutil.NewMockDB(t)

	superUser := dbfs.User{}

	// Getting superuser account
	err := db.Where("email = ?", "superuser").First(&superUser).Error
	assert.NoError(t, err)

	// Creating a file
	rootFolder, err := dbfs.GetRootFolder(db)
	assert.NoError(t, err)
	newFile, err := EXAMPLECreateFile(db, &superUser, "blahblah", rootFolder.FileId)
	assert.NoError(t, err)
	assert.NotNil(t, newFile)

	// Create a few users
	// Creating a user
	user1, err := dbfs.CreateNewUser(db, "permissionCheckUser", "user1Name", 1,
		"permissionCheckUser", "refreshToken", "accessToken", "idToken", "testServer")
	assert.Nil(t, err)
	user2, err := dbfs.CreateNewUser(db, "permissionCheckUser2", "user2Name", 1,
		"permissionCheckUser2", "refreshToken2", "accessToken2", "idToken2", "testServer")
	assert.Nil(t, err)
	user3, err := dbfs.CreateNewUser(db, "permissionCheckUser3", "user3Name", 1,
		"permissionCheckUser3", "refreshToken3", "accessToken3", "idToken3", "testServer")
	assert.Nil(t, err)

	// Putting user3 in a group
	group1, err := dbfs.CreateNewGroup(db, "group1", "balh")
	assert.Nil(t, err)
	err = user3.AddToGroup(db, group1)

	// Create a few files under their user directory
	fileBelongingToUser1, err := EXAMPLECreateFile(db, user1, "fileBelongingToUser1", user1.HomeFolderId)
	assert.NoError(t, err)
	file2BelongingToUser1, err := EXAMPLECreateFile(db, user1, "file2BelongingToUser1", user1.HomeFolderId)
	assert.NoError(t, err)
	fileBelongingToUser2, err := EXAMPLECreateFile(db, user2, "fileBelongingToUser2", user2.HomeFolderId)
	assert.NoError(t, err)
	//fileBelongingToUser3, err := EXAMPLECreateFile(db, user3, "fileBelongingToUser3", user3.HomeFolderId)
	//assert.NoError(t, err)

	t.Run("Play favorites", func(t *testing.T) {

		Assert := assert.New(t)

		// Check that favorites right now is empty

		favorites, err := user1.GetFavoriteFiles(db, 0)
		Assert.NoError(err)
		Assert.Empty(favorites)

		// Add a file to favorites
		err = fileBelongingToUser1.AddToFavorites(db, user1)
		Assert.NoError(err)

		// Check that favorites right now is not empty
		favorites, err = user1.GetFavoriteFiles(db, 0)
		Assert.NoError(err)
		Assert.NotEmpty(favorites)
		Assert.Equal(fileBelongingToUser1.FileId, favorites[0].FileId)

		// Ensuring that another user with another favorite file won't interfer with the first user's favorite file
		err = fileBelongingToUser2.AddToFavorites(db, user2)
		Assert.NoError(err)

		// Check that favorites right now is not empty
		favorites, err = user1.GetFavoriteFiles(db, 0)
		Assert.NoError(err)
		Assert.NotEmpty(favorites)
		Assert.Equal(fileBelongingToUser1.FileId, favorites[0].FileId)

	})

	t.Run("Remove Favorites", func(t *testing.T) {

		Assert := assert.New(t)

		// Trying to remove a file that is not in favorites
		err := file2BelongingToUser1.RemoveFromFavorites(db, user1)
		Assert.NoError(err)

		err = fileBelongingToUser1.RemoveFromFavorites(db, user1)
		Assert.NoError(err)

		// Check that it's removed

		favorites, err := user1.GetFavoriteFiles(db, 0)
		Assert.NoError(err)
		Assert.Empty(favorites)

	})

	t.Run("Share files within internal users", func(t *testing.T) {

		Assert := assert.New(t)
		// First, we'll check that both user2 and user 3 have empty shared to lists.

		sharedTo, err := user2.GetSharedWithUser(db)
		Assert.NoError(err)
		Assert.Empty(sharedTo)

		sharedTo, err = user3.GetSharedWithUser(db)
		Assert.NoError(err)
		Assert.Empty(sharedTo)

		// We are going to share fileBelongingToUser1 to user2
		// and file2BelongingToUser1 to group1, which user3 is a part of.

		// Check that they can't access it first

		_, err = dbfs.GetFileById(db, fileBelongingToUser1.FileId, user2)
		Assert.ErrorIs(err, dbfs.ErrFileNotFound)
		_, err = dbfs.GetFileById(db, file2BelongingToUser1.FileId, user3)
		Assert.ErrorIs(err, dbfs.ErrFileNotFound)

		// Share fileBelongingToUser1 to user2
		err = fileBelongingToUser1.AddPermissionUsers(db, &dbfs.PermissionNeeded{Read: true}, user1, *user2)
		Assert.NoError(err)
		// Share file2BelongingToUser1 to group1
		err = file2BelongingToUser1.AddPermissionGroups(db, &dbfs.PermissionNeeded{Read: true}, user1, *group1)
		Assert.NoError(err)

		var permisisons []dbfs.Permission
		err = db.Where("file_id = ?", file2BelongingToUser1.FileId).Find(&permisisons).Error

		// Check that they can access it now
		_, err = dbfs.GetFileById(db, fileBelongingToUser1.FileId, user2)
		Assert.NoError(err)
		_, err = dbfs.GetFileById(db, file2BelongingToUser1.FileId, user3)
		Assert.NoError(err)

		// Check that their shared links are correct
		sharedTo, err = user2.GetSharedWithUser(db)
		Assert.NoError(err)
		Assert.NotEmpty(sharedTo)
		Assert.Equal(fileBelongingToUser1.FileId, sharedTo[0].FileId)

		sharedTo, err = user3.GetSharedWithUser(db)
		Assert.NoError(err)
		Assert.NotEmpty(sharedTo)
		Assert.Equal(file2BelongingToUser1.FileId, sharedTo[0].FileId)

	})
}
