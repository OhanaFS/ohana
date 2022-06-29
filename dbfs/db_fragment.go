package dbfs

import (
	"gorm.io/gorm"
	"time"
)

const (
	FRAGMENTSTATUSREBUILDING = int8(0) // Fragment is being rebuilt
	FRAGMENTSTATUSGOOD       = int8(1) // Fragment is good
	FRAGMENTSTATUSBAD        = int8(2) // Fragment is bad (hash mismatch), but can be rebuilt
	FRAGMENTSTATUSOFFLINE    = int8(3) // File not found on server
)

type Fragment struct {
	gorm.Model
	FileVersion              FileVersion `gorm:"foreignKey:FileVersionFileId,FileVersionVersionNo"`
	FileVersionFileId        string      `gorm:"primaryKey"`
	FileVersionDataId        string      `gorm:"primaryKey"`
	FileVersionVersionNo     int         `gorm:"primaryKey"`
	FileVersionDataIdVersion int
	FragId                   int `gorm:"primaryKey"`
	ServerId                 string
	FileFragmentPath         string
	Checksum                 string
	LastChecked              time.Time
	FragCount                int `gorm:"not null"`
	Status                   int8
}

// UpdateStatus gets called after checking the status of each fragment
func (f *Fragment) UpdateStatus(tx *gorm.DB, status int8) error {
	f.LastChecked = time.Now()
	f.Status = status
	return tx.Save(f).Error
}

func deleteFragmentsByDataId(tx *gorm.DB, dataId string) error {
	return tx.Where("file_version_data_id = ?", dataId).Delete(&Fragment{}).Error
}
