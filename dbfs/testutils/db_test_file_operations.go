package dbfstestutils

import (
	"github.com/OhanaFS/ohana/dbfs"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	ExampleTotalShards  = 5
	ExampleDataShards   = 3
	ExampleParityShards = 2
	ExampleKeyThreshold = 2
)

type ExampleFile struct {
	FileName       string
	ParentFolderId string
	Server         string
	FragmentPath   string
	FileData       string
	Size           int
	ActualSize     int
}

type ExampleUpdate struct {
	NewSize       int
	NewActualSize int
	Server        string
	Password      string
}

// EXAMPLECreateFile is an example driver for creating a File.
func EXAMPLECreateFile(tx *gorm.DB, user *dbfs.User, fileParams ExampleFile) (*dbfs.File, error) {

	// This is an example script to show how the process should work.

	// First, the system receives the whole file from the user
	// Then, the system creates the record in the system

	file := dbfs.File{
		FileId:             uuid.New().String(),
		FileName:           fileParams.FileName,
		MIMEType:           "",
		ParentFolderFileId: &fileParams.ParentFolderId, // root folder for now
		Size:               fileParams.Size,
		VersioningMode:     dbfs.VersioningOff,
		Checksum:           "CHECKSUM",
		TotalShards:        ExampleTotalShards,
		DataShards:         ExampleDataShards,
		ParityShards:       ExampleParityShards,
		KeyThreshold:       ExampleKeyThreshold,
		PasswordProtected:  false,
		HandledServer:      fileParams.Server,
	}

	fileKey, fileIv, err := dbfs.GenerateKeyIV()
	if err != nil {
		return nil, err
	}

	passwordProtect := dbfs.PasswordProtect{
		FileId:         file.FileId,
		FileKey:        fileKey,
		FileIv:         fileIv,
		PasswordActive: false,
	}

	// This is the key and IV from the pipeline
	dataKey, dataIv, err := dbfs.GenerateKeyIV()
	if err != nil {
		return nil, err
	}

	// Encrypt dataKey and dataIv with the fileKey and fileIv
	dataKey, err = dbfs.EncryptWithKeyIV(dataKey, fileKey, fileIv)
	if err != nil {
		return nil, err
	}
	dataIv, err = dbfs.EncryptWithKeyIV(dataIv, fileKey, fileIv)
	if err != nil {
		return nil, err
	}

	err = tx.Transaction(func(tx *gorm.DB) error {

		err := dbfs.CreateInitialFile(tx, &file, fileKey, fileIv, dataKey, dataIv, user)
		if err != nil {
			return err
		}

		err = tx.Create(&passwordProtect).Error
		if err != nil {
			return err
		}

		err = dbfs.CreatePermissions(tx, &file)
		if err != nil {
			// By right, there should be no error possible? If any error happens, it's likely a system error.
			// However, in the case there is an error, we will revert the transaction (thus deleting the file entry)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// If no error, file is created. Now the system need to process the file in the pipeline and send it to each server
	// Pipeline can get the amount of parity bits based on the Parity Count, and get the amount of shards based on the amount of servers

	// Then, the system splits the files accordingly

	// Then, the system send each shard to each server, and once it's sent successfully

	for i := 1; i <= file.TotalShards; i++ {
		fragId := i
		fragmentPath := uuid.New().String()
		serverId := fileParams.Server

		err = dbfs.CreateFragment(tx, file.FileId, file.DataId, file.VersionNo, fragId, serverId, fragmentPath)
		if err != nil {
			// Not sure how to handle this multiple error situation that is possible.
			// Don't necessarily want to put it in a transaction because I'm worried it'll be too long?
			// or does that make no sense?
			err2 := file.Delete(tx, user, fileParams.Server)
			if err2 != nil {
				return nil, err2
			}

			return nil, err
			// If fails, delete File record and return error.
		}
	}

	// Now we can update it to indicate that it has been saved successfully.

	err = dbfs.FinishFile(tx, &file, user, fileParams.ActualSize, "checksum")
	if err != nil {
		// If fails, delete File record and return error.
	}

	err = dbfs.CreateFileVersionFromFile(tx, &file, user)
	if err != nil {
		// If fails, delete File record and return error.
	}

	return &file, nil

}

func EXAMPLEUpdateFile(tx *gorm.DB, file *dbfs.File, eU ExampleUpdate, user *dbfs.User) error {

	// Key and IV from pipeline
	dataKey, dataIv, err := dbfs.GenerateKeyIV()
	if err != nil {
		return err
	}

	err = file.UpdateFile(tx, eU.NewSize, eU.NewActualSize, "wowKey", eU.Server,
		dataKey, dataIv, eU.Password, user)
	if err != nil {
		return err
	}

	// As each fragment is uploaded, each fragment is added to the database.
	for i := 1; i <= file.TotalShards; i++ {
		err = file.UpdateFragment(tx, i, uuid.NewString(), "wowChecksum", eU.Server)
		if err != nil {
			return err
		}
	}

	// Once all fragments are uploaded, the file is marked as finished.
	return file.FinishUpdateFile(tx, "")

}
