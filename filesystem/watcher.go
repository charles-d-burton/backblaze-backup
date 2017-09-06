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

	"github.com/boltdb/bolt"
	"github.com/fsnotify/fsnotify"
	"github.com/gogo/protobuf/proto"
)

var (
	hashJobs    = make(chan string, 200)
	hashResults = make(chan MetaData, 1000)
	separator   = string(filepath.Separator)
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
	for _, dir := range dirs.Dirs {
		log.Println("Dir to Index: ")
		log.Println(dir)
	}
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
				log.Println("event:", event)

				switch {
				case event.Op&fsnotify.Write == fsnotify.Write ||
					event.Op&fsnotify.Create == fsnotify.Create ||
					event.Op&fsnotify.Rename == fsnotify.Rename:

					accumulator.Lock()
					accumulator.Files[event.Name] = true
					log.Println("Files accumulated: ")
					for accum := range accumulator.Files {
						log.Println(accum)
					}
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
		log.Println("Watching dir: ", dir)
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
				hashJobs <- buffer.String()
				buffer.Reset()
			}
		}
	}
}

//Synchronize ...Synchronize Index with Backblaze
func (dirs *WatchDirs) Synchronize() {
}

func startHashWorkerPool(workers int) {
	for i := 1; i <= workers; i++ {
		go hashWorker(i, hashJobs, hashResults)
		go boltIndexWorker(i, hashResults)
	}
}

func hashWorker(id int, jobs <-chan string, results chan<- MetaData) {
	db := datastores.BoltConn
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(fileIndexName))
		return err
	})
	if err != nil {
		log.Fatal("Unable to open Bolt: ", err)
	}
	for job := range jobs {
		log.Println("Worker: ", id, "started job: ", job)
		file, err := os.Open(job)
		if err != nil {
			log.Println("File error: ", err)
			continue
		}
		stat, err := file.Stat()
		if err != nil {
			log.Println("File error: ", err)
			continue
		}
		err = db.View(func(tx *bolt.Tx) error {
			fileIndexBucket := tx.Bucket([]byte(fileIndexName))
			fileData := fileIndexBucket.Get([]byte(job))
			if fileData != nil {
				fileMetaData := &MetaData{}
				err = proto.Unmarshal(fileData, fileMetaData)
				if err != nil {
					log.Println("Unmarshaling error: ", err)
					return err
				}
				if stat.Size() != fileMetaData.GetSize() {

					hash, err := hashFile(*file)
					if err != nil {
						return err
					}
					fileMetaData := MetaData{
						Name: proto.String(job),
						Size: proto.Int64(stat.Size()),
						Sha1: proto.String(hash),
					}
					results <- fileMetaData
				} else {
					log.Println("File sizes match!")
				}
			} else {
				hash, err := hashFile(*file)
				if err != nil {
					return err
				}
				fileMetaData := MetaData{
					Name: proto.String(job),
					Size: proto.Int64(stat.Size()),
					Sha1: proto.String(hash),
				}
				results <- fileMetaData
			}
			return nil
		})
		//hash, err := hashFile(job)
		if err != nil {
			log.Println(err)
		}
		file.Close()
		//log.Println("Hash: ", hash)
		//results <- HashedFile{FileName: job, Hash: hash}
	}
}

/*
 *
 */
func boltIndexWorker(id int, jobs <-chan MetaData) {
	db := datastores.BoltConn
	var err error
	err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(fileIndexName))
		return err
	})

	for job := range jobs {
		err := db.Update(func(tx *bolt.Tx) error {

			fileIndexBucket := tx.Bucket([]byte(fileIndexName))
			//log.Println("Worker: ", id, "started bolt job: ", job)
			//Get the previously saved object
			fileData := fileIndexBucket.Get([]byte(job.GetName()))
			if fileData != nil {
				fileMetaData := &MetaData{}
				err = proto.Unmarshal(fileData, fileMetaData)
				if err != nil {
					log.Fatal("Unmarshaling error: ", err)
					return err
				}
				//Compare the Sha signatures of the current job with the recorded job
				if job.GetSha1() != fileMetaData.GetSha1() {
					//They don't match so marshal it into a protobuf
					data, err := proto.Marshal(&job)
					if err != nil {
						log.Fatal("Marshaling error: ", err)
						return err
					}
					//Save the new metadata replacing the old object
					fileIndexBucket.Put([]byte(job.GetName()), data)
				}
				//Did not find a file with that name
			} else {
				//Make sure it's not an empty name
				if job.GetName() != "" {
					data, err := proto.Marshal(&job)
					if err != nil {
						log.Fatal("Marshaling error: ", err)
						return err
					}
					fileIndexBucket.Put([]byte(job.GetName()), data)
				}
			}
			return nil
		})
		if err != nil {
			log.Println(err)
		}
	}
}

//Creat a sha1 of a file of any size
func hashFile(file os.File) (string, error) {

	//get information about the file
	info, err := file.Stat()
	if err != nil {
		return "", err
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
		return hex.EncodeToString(hash.Sum(nil)), nil
	}
	return "", err
}
