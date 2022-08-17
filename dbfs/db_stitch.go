package dbfs

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StitchParams struct {
	DataShards   int
	ParityShards int
	KeyThreshold int
}

const (
	defaultDataShards   = 2
	defaultParityShards = 1
	defaultKeyThreshold = 2
)

// SetStitchParams will update stitch parameters as given into dbfs.
// If the value is invalid (i.e. less than 1 or more than 10), it will be set to the default value.
func SetStitchParams(tx *gorm.DB, dataShards, parityShards, keyThreshold int) error {
	// ensure that the number is valid

	errorString := ""

	if dataShards < 1 || dataShards > 10 {
		dataShards = defaultDataShards
		tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&KeyValueDBPair{
			Key:      "dataShards",
			ValueInt: defaultDataShards,
		})
		errorString += fmt.Sprintf("dataShards value :%v is not valid, using default value of %v\n",
			dataShards, defaultDataShards)
	} else {
		tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&KeyValueDBPair{
			Key:      "dataShards",
			ValueInt: dataShards,
		})
	}

	if parityShards < 1 || parityShards > 10 {
		parityShards = defaultParityShards
		tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&KeyValueDBPair{
			Key:      "parityShards",
			ValueInt: defaultParityShards,
		})
		errorString += fmt.Sprintf("parityShards value :%v is not valid, using default value of %v\n",
			parityShards, defaultParityShards)
	} else {
		tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&KeyValueDBPair{
			Key:      "parityShards",
			ValueInt: parityShards,
		})
	}

	if keyThreshold < 1 || keyThreshold > parityShards+dataShards {
		tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&KeyValueDBPair{
			Key:      "keyThreshold",
			ValueInt: defaultKeyThreshold,
		})
		errorString += fmt.Sprintf("keyThreshold value :%v is not valid, using default value of %v\n",
			keyThreshold, defaultKeyThreshold)
	} else {
		tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&KeyValueDBPair{
			Key:      "keyThreshold",
			ValueInt: keyThreshold,
		})
	}
	if errorString != "" {
		return errors.New(errorString)
	} else {
		return nil
	}

}

// GetStitchParams gets the stitch parameters from the dbfs.
// If not found, uses default values (does not write it into DB.)
func GetStitchParams(tx *gorm.DB, log *zap.Logger) (*StitchParams, error) {

	var dataShards, parityShards, keyThreshold KeyValueDBPair
	var dataShardsInt, parityShardsInt, keyThresholdInt int

	// Get dataShards
	err := tx.Where("key = ?", "dataShards").First(&dataShards).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		dataShardsInt = defaultDataShards
		log.Error("dataShards not found, using default value", zap.Int("dataShards", dataShardsInt))
	} else if err != nil {
		return nil, err
	} else {
		dataShardsInt = dataShards.ValueInt
	}

	// Get parityShards
	err = tx.Where("key = ?", "parityShards").First(&parityShards).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		parityShardsInt = defaultParityShards
		log.Error("parityShards not found, using default value", zap.Int("parityShards", parityShardsInt))
	} else if err != nil {
		return nil, err
	} else {
		parityShardsInt = parityShards.ValueInt
	}

	// Get keyThreshold
	err = tx.Where("key = ?", "keyThreshold").First(&keyThreshold).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		keyThresholdInt = defaultKeyThreshold
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
