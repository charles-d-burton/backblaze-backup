package main

import (
	"backblaze-backup/datastores"
	"backblaze-backup/filesystem"
	"log"
	"time"

	"github.com/boltdb/bolt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func main() {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	var err error
	datastores.BoltConn, err = bolt.Open("backblaze.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer datastores.BoltConn.Close()
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/backblaze-backup")
	viper.AddConfigPath("$HOME/.backblaze-backup")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config changed: ", e.Name)
	})
	if err != nil {
		log.Fatal(err)
	}
	//accountId := viper.Get("account-id")
	//applicationKey := viper.Get("application-key")
	watchDirs := viper.GetStringSlice("watch-dirs")
	log.Println(watchDirs)
	filesystem.Watches(watchDirs)
	/*dirs := []string{"/tmp"}

	go filesystem.Watches(dirs)
	for {
	}*/
}
