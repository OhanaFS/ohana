package dbfs

import (
	"errors"
	"github.com/OhanaFS/ohana/util"
	"gorm.io/gorm"
)

func GetToBeDeletedFragments(tx *gorm.DB) ([]Fragment, error) {

	var fragments []Fragment
	var fileVersions []FileVersion

	// Get FileVersions that are marked as to be deleted
	// TODO: Delete needs to have a handled server string.
	err := tx.Where("status = ?", FileStatusToBeDeleted).Find(&fileVersions).Error
	if err != nil {
		return nil, err
	}

	if len(fileVersions) == 0 {
		return fragments, nil
	}

	// Get Fragments that are marked as to be deleted

	dataIdSeen := util.NewSet[string]()

	for _, fileVersion := range fileVersions {

		if dataIdSeen.Has(fileVersion.DataId) {
			continue
		}

		// Check if any of the file version's data fragments are still being used

		var copiesOfData int64
		err2 := tx.Model(&DataCopies{}).Where("data_id = ?", fileVersion.DataId).Count(&copiesOfData).Error
		if err2 != nil {
			return nil, err2
		}
		if copiesOfData >= 1 {
			// There are still data copies of this file
			// Thus, search in the DB to verify that another file is still using it

			var copiesOfDataUsing int64
			err2 = tx.Model(&FileVersion{}).Where("data_id = ? AND status = ?", fileVersion.DataId, FileStatusGood).
				Count(&copiesOfDataUsing).Error

			if copiesOfDataUsing >= 1 {
				continue
			} else {
				// There are no other files using this data fragment
				// Thus, we can delete it.
				err3 := tx.Delete(&DataCopies{}, "data_id = ?", fileVersion.DataId).Error
				if err3 != nil {
					return nil, err3
				}
			}

		}

		// There are no data copies of this file
		// Thus, we can get the fragments and append to fragment.
		var tempFragments []Fragment

		err3 := tx.Where("file_version_data_id = ?", fileVersion.DataId).Find(&tempFragments).Error
		if err3 != nil {
			if errors.Is(err3, gorm.ErrRecordNotFound) {
				continue
			}
			return nil, err3
		}

		fragments = append(fragments, tempFragments...)

		dataIdSeen.Add(fileVersion.DataId)
	}

	return fragments, err
}

func FinishDeleteDataId(tx *gorm.DB, dataId string) error {

	// Delete all fragments with this dataId
	err := tx.Where("file_version_data_id = ?", dataId).Delete(&Fragment{}).Error
	if err != nil {
		return err
	}

	// TODO: Maybe this should just delete. Yeah I don't see a reason to have FileStatusDeleted.
	// In that case, deleting a folder should just go ahead and delete everything else as well.
	return tx.Model(&FileVersion{}).Where("data_id = ?", dataId).Update("status", FileStatusDeleted).Error
}

// MarkOldFileVersions goes through the database and see what file versions are old and can be deleted.
func MarkOldFileVersions(tx *gorm.DB) int {

	return 0
}
