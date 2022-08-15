package dbfs_test

import (
	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/util/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateSharedLinks(t *testing.T) {
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

	// Creating a shared link
	shortenedLink, err := newFile.CreateSharedLink(db, &superUser, "")
	assert.NoError(t, err)
	assert.NotEmpty(t, shortenedLink)

	t.Run("Creating duplicate shortnames", func(t *testing.T) {

		// Creating a shared link
		_, err := newFile.CreateSharedLink(db, &superUser, shortenedLink.ShortenedLink)
		assert.Error(t, err)
		assert.Equal(t, dbfs.ErrLinkExists, err)
	})

	t.Run("Checking that the shortname created is valid", func(t *testing.T) {

		Assert := assert.New(t)
		Assert.NotEmpty(shortenedLink.ShortenedLink)
		t.Log(shortenedLink.ShortenedLink)
		Assert.NotEmpty(shortenedLink.FileId)
		Assert.NotEmpty(shortenedLink.CreatedTime)

		// Checking that we can get the file directly from it.

		retrievedFile, err := dbfs.GetFileFromShortenedLink(db, shortenedLink.ShortenedLink)
		Assert.NoError(err)
		Assert.NotNil(retrievedFile)
		Assert.Equal(newFile.FileId, retrievedFile.FileId)

	})

	t.Run("Modifying the shared link", func(t *testing.T) {

		Assert := assert.New(t)

		// Modifying the shared link
		err := newFile.UpdateSharedLink(db, &superUser, shortenedLink.ShortenedLink, "blah")
		Assert.NoError(err)

		_, err = dbfs.GetFileFromShortenedLink(db, shortenedLink.ShortenedLink)
		Assert.ErrorIs(err, dbfs.ErrSharedLinkNotFound)

		newshortlink, err := dbfs.GetFileFromShortenedLink(db, "blah")
		Assert.NoError(err)
		Assert.Equal(newFile.FileId, newshortlink.FileId)

	})

	t.Run("Deleting the shared link", func(t *testing.T) {

		Assert := assert.New(t)

		// Deleting the shared link
		err := newFile.DeleteSharedLink(db, &superUser, "blah")
		Assert.NoError(err)

		_, err = dbfs.GetFileFromShortenedLink(db, shortenedLink.ShortenedLink)
		Assert.ErrorIs(err, dbfs.ErrSharedLinkNotFound)

		_, err = dbfs.GetFileFromShortenedLink(db, "blah")
		Assert.ErrorIs(err, dbfs.ErrSharedLinkNotFound)

	})

}
