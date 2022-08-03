package dbfs

import (
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type KeyValueDBPair struct {
	Key         string `gorm:"primaryKey"`
	ValueInt    int
	ValueString string
}

type StitchParams struct {
	DataShards   int
	ParityShards int
	KeyThreshold int
}

func GetStitchParams(tx *gorm.DB, log *zap.Logger) (*StitchParams, error) {

	var dataShards, parityShards, keyThreshold KeyValueDBPair
	var dataShardsInt, parityShardsInt, keyThresholdInt int

	// Get dataShards
	err := tx.Where("key = ?", "dataShards").First(&dataShards).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		dataShardsInt = 2
		log.Error("dataShards not found, using default value", zap.Int("dataShards", dataShardsInt))
	} else if err != nil {
		return nil, err
	} else {
		dataShardsInt = dataShards.ValueInt
	}

	// Get parityShards
	err = tx.Where("key = ?", "parityShards").First(&parityShards).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		parityShardsInt = 1
		log.Error("parityShards not found, using default value", zap.Int("parityShards", parityShardsInt))
	} else if err != nil {
		return nil, err
	} else {
		parityShardsInt = parityShards.ValueInt
	}

	// Get keyThreshold
	err = tx.Where("key = ?", "keyThreshold").First(&keyThreshold).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		keyThresholdInt = 2
		log.Error("keyThreshold not found, using default value", zap.Int("keyThreshold", keyThresholdInt))
	} else if err != nil {
		return nil, err
	} else {
		keyThresholdInt = keyThreshold.ValueInt
	}

	return &StitchParams{
		DataShards:   dataShardsInt,
		ParityShards: parityShardsInt,
		KeyThreshold: keyThresholdInt,
	}, nil

}
