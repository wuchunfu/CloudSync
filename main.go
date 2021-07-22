package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/wuchunfu/CloudSync/handler/watchFile"
	"log"
	"os"
	"os/signal"
)

func main() {
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("NewWatcher failed: ", err)
		return
	}

	w := watchFile.Watch{
		Watch: watch,
	}

	defer func(Watch *fsnotify.Watcher) {
		err := Watch.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(w.Watch)

	w.WatchDir("/tmp/test")

	// 终止信号
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, os.Kill)
	<-done
}
