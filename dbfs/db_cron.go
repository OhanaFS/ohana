package dbfs

import (
	"errors"
	"gorm.io/gorm"
)

// Set Type
type dataIdSet map[string]struct{}

// Adds a dataId to the set
func (d dataIdSet) add(dataId string) {
	d[dataId] = struct{}{}
}

// Removes a dataId from the set
func (d dataIdSet) remove(dataId string) {
	delete(d, dataId)
}

// Returns a boolean value describing if the dataId exists in the set
func (d dataIdSet) has(dataId string) bool {
	_, ok := d[dataId]
	return ok
}

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

	dataIdSeen := dataIdSet{}

	for _, fileVersion := range fileVersions {

		if dataIdSeen.has(fileVersion.DataId) {
			continue
		}

		// Check if any of the file version's data fragments are still being used

		var dataCopies DataCopies
		err2 := tx.First(&dataCopies, fileVersion.FileId).Error
		if err2 != nil {
			if !errors.Is(err2, gorm.ErrRecordNotFound) {
				// There are still data copies of this file
				// Thus, search in the DB to verify that another file is still using it

				var dupFileVersions []FileVersion
				err3 := tx.Where("data_id = ? AND status = ?", fileVersion.DataId, FileStatusGood).Find(&dupFileVersions).Error
				if err3 != nil {
					return nil, err3
				}

				if len(dupFileVersions) >= 1 {
					// There are still other files using this data fragment
					// Thus, don't delete it
					continue
				} else {
					// There are no other files using this data fragment
					// Thus, we can delete it.
					tx.Delete(&dataCopies)
				}

			} else {
				return nil, err2
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

		dataIdSeen.add(fileVersion.DataId)
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
