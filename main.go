package main

import (
	"backblaze-backup/database"
	"backblaze-backup/filesystem"
	"log"

	"github.com/boltdb/bolt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
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
