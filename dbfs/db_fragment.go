package dbfs

import (
	"gorm.io/gorm"
	"time"
)

const (
	FragmentStatusRebuilding = int8(0) // Fragment is being rebuilt
	FragmentStatusGood       = int8(1) // Fragment is good
	FragmentStatusBad        = int8(2) // Fragment is bad (hash mismatch), but can be rebuilt
	FragmentStatusOffline    = int8(3) // File not found on server
)

type Fragment struct {
	gorm.Model
	FileVersion              FileVersion `gorm:"foreignKey:FileVersionFileId,FileVersionVersionNo"`
	FileVersionFileId        string      `gorm:"primaryKey"`
	FileVersionDataId        string      `gorm:"primaryKey"`
	FileVersionVersionNo     int         `gorm:"primaryKey"`
	FileVersionDataIdVersion int
	FragId                   int `gorm:"primaryKey"`
	ServerName               string
	FileFragmentPath         string
	LastChecked              time.Time
	TotalShards              int `gorm:"not null"`
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

func GetFragmentByServer(tx *gorm.DB, serverName string) ([]Fragment, error) {
	var fragments []Fragment
	err := tx.Where("server_name = ?", serverName).Find(&fragments).Error
	return fragments, err
}
