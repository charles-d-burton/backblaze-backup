package filesystem

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

var (
	hashJobs    = make(chan string, 200)
	hashResults = make(chan string, 1000)
)

const filechunk = 8192

type WatchDirs struct {
	Dirs []string
}

type HashedFile struct {
	Hash     string
	FileName string
}

//Watches ...Recursively walk the filesystem, entrypoint to file watching
func Watches(tops []string) {
	startHashWorkerPool(8)
	var dirs WatchDirs
	for _, top := range tops {
		err := filepath.Walk(top, func(path string, f os.FileInfo, err error) error {
			if err != nil {
				log.Println(err)
				return err
			}
			//log.Println("File: ", path)
			if f.IsDir() {
				//log.Println("Path: ", path)
				dirs.Dirs = append(dirs.Dirs, path)
			}
			return nil
		})
		if err != nil {
			log.Println(err)
		}
		log.Println("Continuing Loop")
	}
	log.Println("Starting Dedup: ")
	dirs.Dedup()
	dirs.Index()
	log.Println("Post Dedup: ")
	for _, dir := range dirs.Dirs {
		log.Println(dir)
	}
	dirs.Watch()
}

//Dedup ...Remove all duplicate entires from configured directories
func (dirs *WatchDirs) Dedup() {
	log.Println("deduping")
	uniqueSet := make(map[string]bool, len(dirs.Dirs))
	for _, x := range dirs.Dirs {
		uniqueSet[x] = true
	}
	result := make([]string, 0, len(uniqueSet))
	for x := range uniqueSet {
		result = append(result, x)
	}
	dirs.Dirs = result
}

//Watch ...Watch the list of created directories for changes
func (dirs *WatchDirs) Watch() {
	watcher, err := fsnotify.NewWatcher()
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
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
					/* var m runtime.MemStats
					runtime.ReadMemStats(&m)
					log.Printf("\nAlloc = %v\nTotalAlloc = %v\nSys = %v\nNumGC = %v\n\n", m.Alloc/1024, m.TotalAlloc/1024, m.Sys/1024, m.NumGC) */
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

//Index ...Index all the files currently in the system
func (dirs *WatchDirs) Index() {
	for _, dir := range dirs.Dirs {
		files, _ := ioutil.ReadDir(dir)
		for _, file := range files {
			if !file.IsDir() {
				hashJobs <- dir + "/" + file.Name()
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
	}
}

func hashWorker(id int, jobs <-chan string, results chan<- string) {
	for job := range jobs {
		log.Println("Worker: ", id, "started job: ", job)
		hash, err := hashFile(job)
		if err != nil {
			log.Println(err)
			//results <- hash
		}
		log.Println("Hash: ", hash)
	}
}

//Creat a sha1 of a file of unknown size
func hashFile(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer f.Close()

	//get information about the file
	info, err := f.Stat()
	if err != nil {
		return "", err
	}
	size := info.Size()

	blocks := uint64(math.Ceil(float64(size) / float64(filechunk)))
	hash := sha1.New()

	for i := uint64(0); i < blocks; i++ {
		blocksize := int(math.Min(filechunk, float64(size-int64(i*filechunk))))

		buf := make([]byte, blocksize)
		f.Read(buf)
		io.WriteString(hash, string(buf))
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
