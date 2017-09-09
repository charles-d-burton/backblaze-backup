package filesystem

import (
	"backblaze-backup/datastores"
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"sync/atomic"

	"github.com/boltdb/bolt"
	"github.com/fsnotify/fsnotify"
)

var (
	checkFiles               = make(chan string, 2000)
	hasher                   = make(chan string, 20000)
	hashResults              = make(chan *MsgpMetaData, 100)
	separator                = string(filepath.Separator)
	totalIndexedData   int64 = 0
	totalProcessedData int64 = 0
)

const (
	filechunk     = 8192
	fileIndexName = "FileIndex"
)

type WatchDirs struct {
	Dirs []string
}

//Watches ...Recursively walk the filesystem, entrypoint to file watching
func Watches(tops []string) {
	db := datastores.BoltConn
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(fileIndexName))
		return err
	})
	if err != nil {
		log.Println("Unable to open Bolt: ", err)
		return
	}
	startHashWorkerPool(8)
	var dirs WatchDirs
	dirSet := make(map[string]bool)
	for _, top := range tops {
		err := filepath.Walk(top, func(path string, f os.FileInfo, err error) error {
			if err != nil {
				log.Println(err)
				//Return nil because I want to continue processing, if you return something the walker stops
				return nil
			}
			//log.Println("File: ", path)
			if f.IsDir() {
				//Maps can only have one key that matches, duplicates will be overwritten
				dirSet[path] = true
			}
			return nil
		})
		if err != nil {
			log.Println(err)
		}
		//log.Println("Continuing Loop")

	}
	if len(dirSet) > 0 {
		dirSlice := make([]string, 0, len(dirSet))
		for k := range dirSet {
			//log.Println("Processing directory: ", k)
			dirSlice = append(dirSlice, k)
			dirs.Dirs = dirSlice
		}
	}
	/* for _, dir := range dirs.Dirs {
		log.Println("Dir to Index: ")
		log.Println(dir)
	} */
	dirs.Index()

	dirs.Watch()
}

//Watch ...Watch the list of created directories for changes
func (dirs *WatchDirs) Watch() {
	watcher, err := fsnotify.NewWatcher()
	accumulator := datastores.GetAccumulator()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				//log.Println("event:", event)

				switch {
				case event.Op&fsnotify.Write == fsnotify.Write ||
					event.Op&fsnotify.Create == fsnotify.Create ||
					event.Op&fsnotify.Rename == fsnotify.Rename:

					accumulator.Lock()
					if !filterFile(event.Name) {
						accumulator.Files[event.Name] = true
					} else {
						log.Println("Filtering File: ", event.Name)
					}

					/* log.Println("Files accumulated: ")
					for accum := range accumulator.Files {
						log.Println(accum)
					} */
					accumulator.Unlock()
					//hashFile(event.Name)
				default:
					log.Println("no action")
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	for _, dir := range dirs.Dirs {
		//log.Println("Watching dir: ", dir)
		err = watcher.Add(dir)
		if err != nil {
			log.Println(err)
		}
	}
	<-done
}

/*
 *Index ...Index all the files currently in the system
 *Get every file from every configured path and send that
 *to a worker function to be hashed
 */
func (dirs *WatchDirs) Index() {
	var buffer bytes.Buffer
	separator := string(filepath.Separator)
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
				} else {
					log.Println("Filtering File: ", fileName)
				}
				buffer.Reset()
			}
		}
	}
}

func startHashWorkerPool(workers int) {
	for i := 1; i <= workers; i++ {
		go checkFile(i, checkFiles, hasher)
		go boltIndexWorker(i, hashResults)
	}

	//Start only one file worker, this makes the IO on spinning rust a bit better
	go fileWorker(1, hasher, hashResults)

}

func checkFile(id int, jobs <-chan string, results chan<- string) {
	db := datastores.BoltConn
	for job := range jobs {
		file, err := os.Open(job)
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
		sizeMatch := false
		err = db.View(func(tx *bolt.Tx) error {
			fileIndexBucket := tx.Bucket([]byte(fileIndexName))
			fileData := fileIndexBucket.Get([]byte(job))
			if fileData != nil {
				var fileMetaData MsgpMetaData
				_, err := fileMetaData.UnmarshalMsg(fileData)
				if err != nil {
					log.Println("Unmarshaling error: ", err)
					return err
				}
				if stat.Size() != fileMetaData.Size {
					sizeMatch = true
					//results <- job
				} /* else {
					log.Println("File sizes match!")
				} */
			} else {
				sizeMatch = true
				//results <- job
			}
			return nil
		})
		if err != nil {
			log.Println(err)
		}
		if sizeMatch {
			results <- job
		}
	}
}

//Synchronize ...Synchronize Index with Backblaze
func (dirs *WatchDirs) Synchronize() {
}

func fileWorker(id int, jobs <-chan string, results chan<- *MsgpMetaData) {

	for job := range jobs {
		//log.Println("Hasher: ", id, "started job: ", job)

		hash, size, err := hashFile(job)
		if err != nil {
			log.Println(err)
		}
		fileMetaData := MsgpMetaData{
			Name: job,
			Size: size,
			Sha1: hash,
		}
		results <- &fileMetaData
		//hash, err := hashFile(job)
		if err != nil {
			log.Println(err)
		}
		//log.Println("Hash: ", hash)
		//results <- HashedFile{FileName: job, Hash: hash}
	}
}

/*
 *
 */
func boltIndexWorker(id int, jobs <-chan *MsgpMetaData) {
	db := datastores.BoltConn
	for job := range jobs {
		//log.Println("Worker: ", id, "started bolt job: ", job.Name)

		err := db.Batch(func(tx *bolt.Tx) error {

			fileIndexBucket := tx.Bucket([]byte(fileIndexName))
			//log.Println("Worker: ", id, "started bolt job: ", job.GetName())

			//Make sure it's not an empty name
			if job.Name != "" {
				data, err := job.MarshalMsg(nil)
				if err != nil {
					log.Println("Marshaling error: ", err)
					return err
				}
				fileIndexBucket.Put([]byte(job.Name), data)
			}
			return nil
		})
		if err != nil {
			log.Println(err)
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
		size := info.Size()

		blocks := uint64(math.Ceil(float64(size) / float64(filechunk)))
		hash := sha1.New()

		for i := uint64(0); i < blocks; i++ {
			blocksize := int(math.Min(filechunk, float64(size-int64(i*filechunk))))

			buf := make([]byte, blocksize)
			file.Read(buf)
			io.WriteString(hash, string(buf))
		}
		return hex.EncodeToString(hash.Sum(nil)), info.Size(), nil
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
