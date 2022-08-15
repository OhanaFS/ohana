package dbfs

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type SharedLink struct {
	ShortenedLink string    `gorm:"primary_key; unique" json:"shortened_link"`
	FileId        string    `gorm:"primary_key" json:"file_id"`
	CreatedTime   time.Time `json:"created_time"`
}

// GetFileFromShortenedLink Provides the File when given a shortened link
func GetFileFromShortenedLink(db *gorm.DB, shortenedLink string) (*File, error) {
	var sharedLink SharedLink
	err := db.First(&sharedLink, "shortened_link = ?", shortenedLink).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSharedLinkNotFound
		}
		return nil, err
	}

	// Get file from file id
	var file File
	err = db.First(&file, "file_id = ?", sharedLink.FileId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFileNotFound
		}
		return nil, err
	}

	return &file, nil
}
