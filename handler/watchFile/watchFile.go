package watchFile

import (
	"github.com/fsnotify/fsnotify"
	"github.com/wuchunfu/CloudSync/utils"
	"log"
	"os"
	"path/filepath"
)

// NotifyFile 包的指针结构
type NotifyFile struct {
	watch *fsnotify.Watcher
	Path  chan ActionPath
}

// ActionPath 文件操作
type ActionPath struct {
	Path       string
	ActionType fsnotify.Op
	desc       string
	SourcePath string
	TargetPath string
}

// NewNotifyFile 返回 fsnotify 对象指针
func NewNotifyFile() *NotifyFile {
	notifyFile := new(NotifyFile)
	notifyFile.watch, _ = fsnotify.NewWatcher()
	notifyFile.Path = make(chan ActionPath, 10)
	return notifyFile
}

// WatchDir a directory
func (notifyFile *NotifyFile) WatchDir(sourcePath string, targetPath string) {
	fullPath := utils.GetFullPath(sourcePath)
	log.Println("Watching:", fullPath)
	// Walk all directory
	err := filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
		// 判断是否为目录, 监控目录以及目录下文件, 目录下的文件也在监控范围内
		// Just watch directory(all child can be watched)
		if info.IsDir() {
			err = notifyFile.watch.Add(path)
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

	// 协程
	go notifyFile.WatchEvents(fullPath, targetPath)
}

// WatchEvents Handle the watch events 监控目录
func (notifyFile *NotifyFile) WatchEvents(sourcePath string, targetPath string) {
	for {
		select {
		case event := <-notifyFile.watch.Events:
			//log.Println("event:", event)
			// Create event
			if event.Op&fsnotify.Create == fsnotify.Create {
				if !utils.IgnoreFile(event.Name) {
					log.Println(">>", event.Name, "[create]")
					go notifyFile.PushEventChannel(event.Name, fsnotify.Create, "添加监控", sourcePath, targetPath)
				}
				// 获取新创建文件的信息, 如果是目录, 则加入监控中
				if utils.IsDir(event.Name) {
					err := notifyFile.watch.Add(event.Name)
					if err != nil {
						log.Println(err)
					}
					log.Println("添加监控: ", event.Name)
				}
			}

			// write event
			if event.Op&fsnotify.Write == fsnotify.Write {
				if !utils.IgnoreFile(event.Name) {
					log.Println(">>", event.Name, "[edit]")
					go notifyFile.PushEventChannel(event.Name, fsnotify.Write, "写入文件", sourcePath, targetPath)
				}
			}

			// delete event
			if event.Op&fsnotify.Remove == fsnotify.Remove {
				// 如果删除文件是目录，则移除监控
				if utils.IsDir(event.Name) {
					err := notifyFile.watch.Remove(event.Name)
					if err != nil {
						log.Println(err)
					}
					log.Println("删除监控: ", event.Name)
				}
				if !utils.IgnoreFile(event.Name) {
					log.Println(">>", event.Name, "[remove]")
				}
			}

			// Rename
			if event.Op&fsnotify.Rename == fsnotify.Rename {
				// 如果重命名文件是目录，则移除监控
				// 注意这里无法使用os.Stat来判断是否是目录了
				// 因为重命名后，go已经无法找到原文件来获取信息了
				// 所以这里就简单粗爆的直接remove好了
				err := notifyFile.watch.Remove(event.Name)
				if err != nil {
					log.Println(err)
				}
				if !utils.IgnoreFile(event.Name) {
					log.Println(">>", event.Name, "[rename]")
				}
			}
			// Chmod
			if event.Op&fsnotify.Chmod == fsnotify.Chmod {
				if !utils.IgnoreFile(event.Name) {
					log.Println(">>", event.Name, "[chmod]")
					go notifyFile.PushEventChannel(event.Name, fsnotify.Chmod, "修改权限", sourcePath, targetPath)
				}
			}
		case err := <-notifyFile.watch.Errors:
			log.Println(err)
			return
		}
	}
}

// PushEventChannel 将发生事件加入 channel
func (notifyFile *NotifyFile) PushEventChannel(Path string, ActionType fsnotify.Op, desc string, source string, target string) {
	notifyFile.Path <- ActionPath{
		Path:       Path,
		ActionType: ActionType,
		desc:       desc,
		SourcePath: source,
		TargetPath: target,
	}
}
