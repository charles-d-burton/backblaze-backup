package filesystem

import (
	"backblaze-backup/datastores"
	"log"

	"github.com/boltdb/bolt"
)

const (
	fileIndexName = "FileIndex"
	lastFullIndex = "LastFullIndex"
)

func initBolt() {
	db := datastores.BoltConn
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(fileIndexName))
		_, err = tx.CreateBucketIfNotExists([]byte(lastFullIndex))
		return err
	})
	if err != nil {
		log.Println("Unable to open Bolt: ", err)
		return
	}
}
