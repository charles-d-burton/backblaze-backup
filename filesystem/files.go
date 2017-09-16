package filesystem

import (
	"backblaze-backup/datastores"
	"log"

	"github.com/boltdb/bolt"
	"github.com/gogo/protobuf/proto"
)

//go:generate protoc --go_out=. files.proto

//RecordMetaData ... Take a MetaData object and record it in boltDB
func (md *MetaData) RecordMetaData() error {
	db := datastores.BoltConn
	//log.Println("Worker: ", id, "started bolt job: ", job.Name)

	err := db.Batch(func(tx *bolt.Tx) error {

		fileIndexBucket := tx.Bucket([]byte(fileIndexName))
		//log.Println("Worker: ", id, "started bolt job: ", job.GetName())

		//Make sure it's not an empty name
		if md.Name != "" {
			data, err := proto.Marshal(md)
			if err != nil {
				log.Println("Marshaling error: ", err)
				return err
			}
			fileIndexBucket.Put([]byte(md.Name), data)
		}
		return nil
	})
	return err
}
