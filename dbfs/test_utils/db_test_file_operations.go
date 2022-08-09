package dbfstestutils

import (
	"bytes"
	"encoding/hex"
	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/stitch"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"io"
	"os"
	"path"
	"strconv"
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
	FragmentPath  string
	FileData      string
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

	encoder := stitch.NewEncoder(&stitch.EncoderOptions{
		DataShards:   ExampleDataShards,
		ParityShards: ExampleParityShards,
		KeyThreshold: ExampleKeyThreshold,
	})

	// Creating the temp files
	shardWriters := make([]io.Writer, file.TotalShards)
	shardFiles := make([]*os.File, file.TotalShards)
	shardNames := make([]string, file.TotalShards)

	for i := 1; i <= file.TotalShards; i++ {
		shardName := file.DataId + ".shard" + strconv.Itoa(i)
		shardPath := path.Join(fileParams.FragmentPath, shardName)
		shardFile, err := os.Create(shardPath)
		if err != nil {
			return nil, err
		}
		shardFiles[i-1] = shardFile
		shardWriters[i-1] = shardFile
		shardNames[i-1] = shardName
		defer shardFile.Close()
	}

	// Passing it into encoder

	// Creating a fake file with multipart streamer
	dataAsBytes := []byte(fileParams.FileData)
	dataReader := bytes.NewReader(dataAsBytes)

	dataKeyBytes, err := hex.DecodeString(dataKey)
	if err != nil {
		return nil, err
	}
	dataIvBytes, err := hex.DecodeString(dataIv)
	if err != nil {
		return nil, err
	}

	// Encoding the data
	encode, err := encoder.Encode(dataReader, shardWriters, dataKeyBytes, dataIvBytes)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	for i := 0; i < file.TotalShards; i++ {
		if err = encoder.FinalizeHeader(shardFiles[i]); err != nil {
			return nil, err
		}
	}

	for i := 1; i <= file.TotalShards; i++ {
		fragId := i
		serverId := fileParams.Server
		shardName := shardNames[i-1]

		err = dbfs.CreateFragment(tx, file.FileId, file.DataId, file.VersionNo, fragId, serverId, shardName)
		if err != nil {

			err2 := file.Delete(tx, user, fileParams.Server)
			if err2 != nil {
				return nil, err2
			}

			return nil, err
			// If fails, delete File record and return error.
		}
	}

	// Now we can update it to indicate that it has been saved successfully.

	checksum := hex.EncodeToString(encode.FileHash)

	err = dbfs.FinishFile(tx, &file, user, fileParams.ActualSize, checksum)
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

	encoder := stitch.NewEncoder(&stitch.EncoderOptions{
		DataShards:   ExampleDataShards,
		ParityShards: ExampleParityShards,
		KeyThreshold: ExampleKeyThreshold,
	})

	// Key and IV from pipeline
	dataKey, dataIv, err := dbfs.GenerateKeyIV()
	if err != nil {
		return err
	}

	err = file.UpdateFile(tx, eU.NewSize, eU.NewActualSize, "UPDATING",
		eU.Server, dataKey, dataIv, eU.Password, user)
	if err != nil {
		return err
	}

	shardWriters := make([]io.Writer, ExampleTotalShards)
	shardFiles := make([]*os.File, ExampleTotalShards)
	shardNames := make([]string, ExampleTotalShards)
	for i := 0; i < ExampleTotalShards; i++ {
		shardNames[i] = file.DataId + ".shard" + strconv.Itoa(i+1)
		shardPath := path.Join(eU.FragmentPath, shardNames[i])
		shardFile, err := os.Create(shardPath)
		if err != nil {
			return err
		}
		shardFiles[i] = shardFile
		shardWriters[i] = shardFile
		defer shardFile.Close()
	}

	// encoder
	data := []byte(eU.FileData)
	dataReader := bytes.NewReader(data)
	dataKeyBytes, err := hex.DecodeString(dataKey)
	if err != nil {
		return err
	}
	dataIvBytes, err := hex.DecodeString(dataIv)
	if err != nil {
		return err
	}

	result, err := encoder.Encode(dataReader, shardWriters, dataKeyBytes, dataIvBytes)
	if err != nil {
		return err
	}

	for i := 0; i < ExampleTotalShards; i++ {
		if err = encoder.FinalizeHeader(shardFiles[i]); err != nil {
			return err
		}
	}

	// As each fragment is uploaded, each fragment is added to the database.
	for i := 0; i < ExampleTotalShards; i++ {
		fragId := i + 1

		err = file.UpdateFragment(tx, fragId, shardNames[i], "wowChecksum", eU.Server)
		if err != nil {
			return err
		}
	}

	// Once all fragments are uploaded, the file is marked as finished.
	return file.FinishUpdateFile(tx, hex.EncodeToString(result.FileHash))

}

func EXAMPLECorruptFragments(fragmentPath string) error {

	// open shard
	shardFile, err := os.OpenFile(fragmentPath, os.O_RDWR, 0666)
	if err != nil {
		return err
	}

	_, err = shardFile.Seek(1024, io.SeekStart)
	if err != nil {
		return err
	}

	// damage

	_, err = shardFile.Write([]byte("corruption"))
	if err != nil {
		return err
	}

	err = shardFile.Close()
	if err != nil {
		return err
	}

	return nil
}
