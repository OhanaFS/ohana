package dbfs

// DataCopies This table is used for keeping track of fragments that have more than 1 copy
// and thus can't be deleted without ensuring that that data fragment is no longer being used.
type DataCopies struct {
	DataId string `gorm:"primaryKey; not null; unique"`
}
