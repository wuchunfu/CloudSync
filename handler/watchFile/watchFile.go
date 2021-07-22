package watchFile

import (
	"github.com/fsnotify/fsnotify"
	"github.com/wuchunfu/CloudSync/utils"
	"log"
	"os"
	"path/filepath"
)

var (
	// Global chan variables
	// file_watcher will write the chan and file_handle will read the chan
	// create file
	fileCreateEvent = make(chan string)

	// write
	fileWriteEvent = make(chan string)

	// remove
	fileRemoveEvent = make(chan string)

	// rename
	fileRenameEvent = make(chan string)

	// chmod
	fileChmodEvent = make(chan string)
)

type Watch struct {
	Watch *fsnotify.Watcher
}

// WatchDir a directory
func (w *Watch) WatchDir(dir string) {
	fullPath := utils.GetFullPath(dir)
	log.Println("Watching:", fullPath)
	// Walk all directory
	err := filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
		// 判断是否为目录, 监控目录以及目录下文件, 目录下的文件也在监控范围内
		// Just watch directory(all child can be watched)
		if info.IsDir() {
			err = w.Watch.Add(path)
			if err != nil {
				log.Fatal(err)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	go EventsHandler(w)
}

// Handle the watch events
func EventsHandler(w *Watch) {
	for {
		select {
		case event := <-w.Watch.Events:
			//log.Println("event:", event)
			// Create event
			if event.Op&fsnotify.Create == fsnotify.Create {
				if !utils.IgnoreFile(event.Name) {
					log.Println(">>", event.Name, "[create]")
					fileCreateEvent <- event.Name
				}
				if utils.IsDir(event.Name) {
					err := w.Watch.Add(event.Name)
					if err != nil {
						log.Fatal(err)
					}
				}
			}

			// write event
			if event.Op&fsnotify.Write == fsnotify.Write {
				if !utils.IgnoreFile(event.Name) {
					log.Println(">>", event.Name, "[edit]")
					fileWriteEvent <- event.Name
				}
			}

			// delete event
			if event.Op&fsnotify.Remove == fsnotify.Remove {
				if utils.IsDir(event.Name) {
					err := w.Watch.Remove(event.Name)
					if err != nil {
						log.Fatal(err)
					}
				}
				if !utils.IgnoreFile(event.Name) {
					log.Println(">>", event.Name, "[remove]")
					fileRemoveEvent <- event.Name
				}
			}

			// Rename
			if event.Op&fsnotify.Rename == fsnotify.Rename {
				err := w.Watch.Remove(event.Name)
				if err != nil {
					log.Fatal(err)
				}
				if !utils.IgnoreFile(event.Name) {
					log.Println(">>", event.Name, "[rename]")
					fileRenameEvent <- event.Name
				}
			}
			// Chmod
			if event.Op&fsnotify.Chmod == fsnotify.Chmod {
				if !utils.IgnoreFile(event.Name) {
					log.Println(">>", event.Name, "[chmod]")
					fileChmodEvent <- event.Name
				}
			}
		case err := <-w.Watch.Errors:
			log.Fatal(err)
			return
		}
	}
}
