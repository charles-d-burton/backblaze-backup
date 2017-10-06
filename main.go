package main

import (
	"backblaze-backup/datastores"
	"backblaze-backup/filesystem"
	"backblaze-backup/requests"
	"log"
	"runtime"
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
	globals := viper.GetStringSlice("ignore-filters.global")
	if runtime.GOOS == "linux" {
		filters := viper.GetStringSlice("ignore-filters.linux")
		filters = append(filters, globals...)
		datastores.SetFilters(filters)
	} else if runtime.GOOS == "windows" {
		filters := viper.GetStringSlice("ignore-filters.windows")
		filters = append(filters, globals...)
		datastores.SetFilters(filters)
	} else if runtime.GOOS == "darwin" {
		filters := viper.GetStringSlice("ignore-filters.apple")
		filters = append(filters, globals...)
		datastores.SetFilters(filters)
	}
	log.Println(watchDirs)
	initializeBackblaze()
	filesystem.Watches(watchDirs)
}

func initializeBackblaze() {
	authorization, err := requests.GetAuthorization(viper.GetString("account-id"), viper.GetString("application-key"))
	if err != nil {
		log.Println(err)
		return
	}
	authorization.CreateBackblazeBucket()
}
