package main

import (
	"github.com/wuchunfu/CloudSync/config"
	"github.com/wuchunfu/CloudSync/handler/watchFile"
	"github.com/wuchunfu/CloudSync/utils/sftpUtils"
	"log"
	"os"
	"os/signal"
)

func main() {

	go func() {
		ch := make(chan os.Signal)
		// 获取程序退出信号
		signal.Notify(ch, os.Interrupt, os.Kill)
		<-ch
		log.Println("server exit")
		os.Exit(1)
	}()

	watch := watchFile.NewNotifyFile()
	for _, v := range config.GlobalObject.Sync {
		// 添加监控目录
		watch.WatchDir(v.SourcePath, v.TargetPath)
	}

	sftpClient := sftpUtils.NewSftpHandler()

	go func(*watchFile.NotifyFile) {
		for {
			select {
			case path := <-watch.Path:
				sftpClient.Upload(path.Path, path.TargetPath)
			default:
				continue
			}
		}
	}(watch)

	select {}
}
