package filesystem

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type WatchDirs struct {
	Dirs []string
}

//Watches ...Recursively walk the filesystem
func Watches(tops []string) error {
	var dirs WatchDirs
	for _, top := range tops {
		filepath.Walk(top, func(path string, f os.FileInfo, err error) error {
			if err != nil {
				log.Println(err)
				return nil
			}
			log.Println("Path: ", path)
			dirs.Dirs = append(dirs.Dirs, path)
			return nil
		})
		log.Println("Continuing Loop")
	}
	log.Println("Staring Dedup: ")
	dirs.Dedup()
	log.Println("Post Dedup: ")
	for _, dir := range dirs.Dirs {
		log.Println(dir)
	}
	dirs.Watch()
	return nil
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
}

//Synchronize ...Synchronize Index with Backblaze
func (dirs *WatchDirs) Synchronize() {
}
