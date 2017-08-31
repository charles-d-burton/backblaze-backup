package main

import (
	"backblaze-backup/database"
	"backblaze-backup/filesystem"
	"log"

	"github.com/boltdb/bolt"
)

func main() {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	var err error
	database.BoltConn, err = bolt.Open("backblaze.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer database.BoltConn.Close()
	dirs := []string{"/tmp"}
	go filesystem.Watches(dirs)
	for {
	}
}
