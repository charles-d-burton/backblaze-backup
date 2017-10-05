package filesystem

import (
	"backblaze-backup/datastores"
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gogo/protobuf/proto"
)

var (
	checkFiles = make(chan string, 2000)
	hasher     = make(chan string, 20000)
	separator  = string(filepath.Separator)

	totalIndexedData   int64
	totalProcessedData int64

	fullIndex = false
)

/*
 *Index ...Index all the files currently in the system
 *Get every file from every configured path and send that
 *to a worker function to be hashed
 */
func (dirs *WatchDirs) Index(full bool) {
	dirs.mutex.Lock()
	defer dirs.mutex.Unlock()
	var buffer bytes.Buffer
	separator := string(filepath.Separator)
	fullIndex = full
	for _, dir := range dirs.Dirs {
		files, _ := ioutil.ReadDir(dir)
		for _, file := range files {
			if !file.IsDir() {
				buffer.WriteString(dir)
				buffer.WriteString(separator)
				buffer.WriteString(file.Name())
				fileName := buffer.String()
				if !filterFile(fileName) {
					checkFiles <- fileName
				}
				buffer.Reset()
			}
		}
	}
	//Reset the indexer when indexing is complete
	fullIndex = false
}

//ScheduleIndex ... Schedule a full Reindex of the database
func (dirs *WatchDirs) ScheduleIndex() {
	//Schedule the check for every 30 seconds
	ticker := time.NewTicker(time.Second * 30)
	for t := range ticker.C {
		log.Println("Checking for Reindex", t)
		db := datastores.BoltConn
		err := db.Update(func(tx *bolt.Tx) error {
			lastIndexBucket := tx.Bucket([]byte(lastFullIndexBucket))
			last := lastIndexBucket.Get([]byte("last"))
			//No last index done, schedule it starting now
			if last == nil {
				log.Println("No scheduled index ran")
				b := make([]byte, 8)
				log.Println("Updating with time: ", time.Now().Unix())
				binary.LittleEndian.PutUint64(b, uint64(time.Now().Unix()))
				lastIndexBucket.Put([]byte("last"), b)
			} else {
				//Check if time difference is greater than 30 days (2592000 seconds)
				if (uint64(time.Now().Unix()) - uint64(binary.LittleEndian.Uint64(last))) > 2592000 {
					log.Println("Time since last reindex is greater than threshold!  Running a full index!")
					b := make([]byte, 8)
					binary.LittleEndian.PutUint64(b, uint64(time.Now().Unix()))
					lastIndexBucket.Put([]byte("last"), b)
					go dirs.Index(true)
				}
			}
			return nil
		})
		if err != nil {
			log.Println(err)
		}
	}
}

func startHashWorkerPool(workers int) {
	for i := 1; i <= workers; i++ {
		go checkFile(i, checkFiles, hasher)

	}
	go fileWorker(1, hasher)
	//Start only one file worker, this makes the IO on spinning rust a bit better
}

func checkFile(id int, jobs <-chan string, results chan<- string) {

	for fileName := range jobs {
		file, err := os.Open(fileName)
		if err != nil {
			log.Println("File error: ", err)
			file.Close()
			continue
		}
		stat, err := file.Stat()
		if err != nil {
			log.Println("File error: ", err)
			file.Close()
			continue
		}
		file.Close()
		atomic.AddInt64(&totalProcessedData, stat.Size())
		//log.Println("Total Data Processed: ", bytefmt.ByteSize(uint64(totalProcessedData)))
		//Optionally bypass the size check and force a full Reindex
		if !fullIndex {
			sizeMatch, err := checkFileBySize(fileName, stat.Size())
			if err != nil {
				log.Println(err)
			} else {
				if !sizeMatch {
					results <- fileName
				}
			}
		} else {
			log.Println("Forcing a full reindex")
			results <- fileName
		}
	}
}

/*
 *Compare the file size against the size in the index.  If they match
 *it probably didn't change.
 */
func checkFileBySize(name string, size int64) (bool, error) {
	db := datastores.BoltConn
	match := true
	err := db.View(func(tx *bolt.Tx) error {
		fileIndexBucket := tx.Bucket([]byte(fileIndexName))
		fileData := fileIndexBucket.Get([]byte(name))
		if fileData != nil {
			var fileMetaData MetaData
			err := proto.Unmarshal(fileData, &fileMetaData)
			if err != nil {
				log.Println("Unmarshaling error: ", err)
				return err
			}
			if size != fileMetaData.Size {
				match = false
			}
		} else {
			match = false
		}
		return nil
	})
	return match, err
}

func fileWorker(id int, jobs <-chan string) {

	for job := range jobs {
		//log.Println("Hasher: ", id, "started job: ", job)

		hash, size, err := hashFile(job)
		if err != nil {
			log.Println(err)
		}
		fileMetaData := MetaData{
			Name:     job,
			Size:     size,
			Sha1:     hash,
			BackedUp: false,
		}
		err = fileMetaData.RecordMetaData()
		if err != nil {
			//TODO:  Possibly attempt to index again with a backoff
			log.Println("error Recording Index: ", err)
		}
	}
}

//Creat a sha1 of a file of any size
func hashFile(f string) (string, int64, error) {
	file, err := os.Open(f)
	if err != nil {
		return "", 0, err
	}
	defer file.Close()
	//get information about the file
	info, err := file.Stat()
	if err != nil {
		return "", 0, err
	}
	if !info.IsDir() {
		hasher := sha1.New()
		if _, err := io.Copy(hasher, file); err != nil {
			log.Println("Hashing Error: ", err)
		}
		return hex.EncodeToString(hasher.Sum(nil)), info.Size(), nil
	}
	return "", 0, err
}

func filterFile(file string) bool {
	for _, filter := range datastores.GetFilters() {

		matched := filter.MatchString(file)
		if matched {
			//log.Println("Filtering on: ", filter.String())
			//log.Println("Filtering file: ", file)
			return matched
		}
	}
	return false
}
