package dbfs

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ! To be implemented

type Server struct {
	gorm.Model
	Name string `gorm:"not null"`
}

type Role struct {
	RoleID      int `gorm:"primaryKey;autoIncrement"`
	RoleMapping string
	RoleName    string   `gorm:"not null"`
	Users       []*User  `gorm:"many2many:user_roles;"`
	Groups      []*Group `gorm:"many2many:group_roles;"`
}

type StitchParams struct {
	Key         string `gorm:"primaryKey"`
	ValueInt    int
	ValueString string
}

func GetStitchParams(tx *gorm.DB, log *zap.Logger) (int, int, int, error) {

	var dataShards, parityShards, keyThreshold StitchParams
	var dataShardsInt, parityShardsInt, keyThresholdInt int

	// Get dataShards
	err := tx.Where("key = ?", "dataShards").First(&dataShards).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		dataShardsInt = 2
		log.Error("dataShards not found, using default value", zap.Int("dataShards", dataShardsInt))
	} else if err != nil {
		return 0, 0, 0, err
	} else {
		dataShardsInt = dataShards.ValueInt
	}

	// Get parityShards
	err = tx.Where("key = ?", "parityShards").First(&parityShards).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		parityShardsInt = 1
		log.Error("parityShards not found, using default value", zap.Int("parityShards", parityShardsInt))
	} else if err != nil {
		return 0, 0, 0, err
	} else {
		parityShardsInt = parityShards.ValueInt
	}

	// Get keyThreshold
	err = tx.Where("key = ?", "keyThreshold").First(&keyThreshold).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		keyThresholdInt = 2
		log.Error("keyThreshold not found, using default value", zap.Int("keyThreshold", keyThresholdInt))
	} else if err != nil {
		return 0, 0, 0, err
	} else {
		keyThresholdInt = keyThreshold.ValueInt
	}

	return dataShardsInt, parityShardsInt, keyThresholdInt, nil

}
