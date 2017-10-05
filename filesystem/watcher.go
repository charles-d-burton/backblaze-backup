package filesystem

import (
	"backblaze-backup/datastores"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

type WatchDirs struct {
	mutex sync.RWMutex
	Dirs  []string
}

/*
 *Watches ...Recursively walk the filesystem, entrypoint to file watching
 *Generates a map of directories to watch then converst that map to a slice
 *the slice is used to index files
 */
func Watches(tops []string) {
	initBolt()
	startHashWorkerPool(8)
	var dirs WatchDirs
	dirs.mutex.Lock()
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
	}
	if len(dirSet) > 0 {
		dirSlice := make([]string, 0, len(dirSet))
		for k := range dirSet {
			//log.Println("Processing directory: ", k)
			dirSlice = append(dirSlice, k)
			dirs.Dirs = dirSlice
		}
	}
	dirs.mutex.Unlock()
	go dirs.Index(false)
	go dirs.Watch()
	go dirs.ScheduleIndex()
	//Block forever
	select {}
}

//Watch ...Watch the list of created directories for changes
func (dirs *WatchDirs) Watch() {
	watcher, err := fsnotify.NewWatcher()
	accumulator := datastores.GetAccumulator()
	go scheduleSweep()
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
					}

					/* log.Println("Files accumulated: ")
					for accum := range accumulator.Files {
						log.Println(accum)
					} */
					accumulator.Unlock()
				//hashFile(event.Name)
				case event.Op&fsnotify.Remove == fsnotify.Remove:
					accumulator.Lock()
					delete(accumulator.Files, event.Name)
					accumulator.Unlock()
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

//Clear out the accumulated files every five minutes.
func scheduleSweep() {
	ticker := time.NewTicker(time.Minute * 5)
	for t := range ticker.C {
		accumulator := datastores.GetAccumulator()
		accumulator.Lock()
		log.Println("Clearing accumulated file changes: ", t)
		for file := range accumulator.Files {
			log.Println("Clearing file: ", file)
			checkFiles <- file
			delete(accumulator.Files, file)
		}
		accumulator.Unlock()

	}

}

//Synchronize ...Synchronize Index with Backblaze
func (dirs *WatchDirs) Synchronize() {
}
