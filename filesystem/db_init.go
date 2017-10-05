package filesystem

import (
	"backblaze-backup/datastores"

	"github.com/boltdb/bolt"
)

const (
	fileIndexName       = "FileIndex"
	lastFullIndexBucket = "LastFullIndex"
)

func initBolt() {
	db := datastores.BoltConn
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(fileIndexName))
		_, err = tx.CreateBucketIfNotExists([]byte(lastFullIndexBucket))
		return err
	})
	if err != nil {
		panic(err)
	}
}
