package api

import (
	"container/list"
	"github.com/wuchunfu/CloudSync/common"
	"github.com/wuchunfu/CloudSync/handler/watchx"
	"github.com/wuchunfu/CloudSync/middleware/configx"
	"github.com/wuchunfu/CloudSync/utils/sftpx"
	"log"
	"os"
	"os/signal"
	"runtime"
)

func Run() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	common.Md5Map = make(map[string]string)
	common.WatcherMap = make(map[string]bool) // 监听的文件夹列表
	common.ChangedMap = make(map[int]*list.List)

	go func() {
		ch := make(chan os.Signal)
		// 获取程序退出信号
		signal.Notify(ch, os.Interrupt, os.Kill)
		<-ch
		log.Println("server exit")
		os.Exit(1)
	}()

	watch := watchx.NewNotifyFile()
	for _, v := range configx.ServerSetting.Sync {
		// 添加监控目录
		watch.WatchDir(v.SourcePath, v.TargetPath)
	}

	sftpClient := sftpx.NewSftpHandler()

	go func(*watchx.NotifyFile) {
		for {
			select {
			case path := <-watch.Path:
				sftpClient.Upload(path.Path, path.TargetPath)
			default:
				continue
			}
		}
	}(watch)

	// 重新加载所有MD5,生成新的的csv文件中
	watchx.OutPutToFile()
	log.Println("load scv file done!")

	// 定时任务
	go watchx.TimerCheck()

	select {}
}
