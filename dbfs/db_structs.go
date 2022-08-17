package dbfs

type KeyValueDBPair struct {
	Key         string `gorm:"primaryKey"`
	ValueInt    int
	ValueString string
}
