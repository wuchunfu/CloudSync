package watchFile

import (
	"container/list"
	"encoding/csv"
	"github.com/fsnotify/fsnotify"
	"github.com/wuchunfu/CloudSync/common"
	"github.com/wuchunfu/CloudSync/utils"
	"log"
	"os"
	"path/filepath"
	"time"
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
		if info.IsDir() { // 路径为目录
			notifyFile.LetItWatcher(path)
		} else { // 路径为文件
			name := info.Name()
			if common.OutputFileName == name {
				log.Fatal("不能监控输出文件...")
			}
		}
		// 生成 MD5 值
		if md5Str := utils.GenerateMd5(path); md5Str != "" {
			common.Md5Map[path] = md5Str
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
					// 获取新创建文件的信息, 如果是目录, 则加入监控中
					if utils.IsDir(event.Name) {
						notifyFile.LetItWatcher(event.Name)
					}
					newMd5Str := utils.GenerateMd5(event.Name)
					if len(newMd5Str) > 0 {
						common.Md5Map[event.Name] = newMd5Str
						LetItChanged(1, event.Name)
					}
					go notifyFile.PushEventChannel(event.Name, fsnotify.Create, "添加监控", sourcePath, targetPath)
				}
			}

			// write event
			if event.Op&fsnotify.Write == fsnotify.Write {
				if !utils.IgnoreFile(event.Name) {
					log.Println(">>", event.Name, "[edit]")
					// 判断文件是否存在
					if oldMd5, ok := common.Md5Map[event.Name]; ok {
						// 判断修改前和修改后的 md5 是否一致
						newMd5Str := utils.GenerateMd5(event.Name)
						if oldMd5 != newMd5Str {
							common.Md5Map[event.Name] = newMd5Str
							LetItChanged(1, event.Name)
						}
					} else {
						common.Md5Map[event.Name] = utils.GenerateMd5(event.Name)
						LetItChanged(1, event.Name)
					}
					go notifyFile.PushEventChannel(event.Name, fsnotify.Write, "写入文件", sourcePath, targetPath)
				}
			}

			// delete event
			if event.Op&fsnotify.Remove == fsnotify.Remove {
				if !utils.IgnoreFile(event.Name) {
					log.Println(">>", event.Name, "[remove]")
					if _, ok := common.Md5Map[event.Name]; ok {
						delete(common.Md5Map, event.Name)
						LetItChanged(3, event.Name)
					}
				}
			}

			// Rename
			if event.Op&fsnotify.Rename == fsnotify.Rename {
				if !utils.IgnoreFile(event.Name) {
					log.Println(">>", event.Name, "[rename]")
					if _, ok := common.Md5Map[event.Name]; ok {
						delete(common.Md5Map, event.Name)
						LetItChanged(4, event.Name)
					}
					notifyFile.DeleteItWatcher(event.Name)
				}
			}
			// Chmod
			if event.Op&fsnotify.Chmod == fsnotify.Chmod {
				if !utils.IgnoreFile(event.Name) {
					log.Println(">>", event.Name, "[chmod]")
					if _, ok := common.Md5Map[event.Name]; ok {
						LetItChanged(5, event.Name)
					}
					go notifyFile.PushEventChannel(event.Name, fsnotify.Chmod, "修改权限", sourcePath, targetPath)
				}
			}
		case err := <-notifyFile.watch.Errors:
			log.Println("======", err)
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

// LetItWatcher 加入监控集合中
func (notifyFile *NotifyFile) LetItWatcher(path string) {
	if _, ok := common.WatcherMap[path]; !ok {
		common.WatcherMap[path] = true
		err := notifyFile.watch.Add(path)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// DeleteItWatcher 从监控集合中删除
func (notifyFile *NotifyFile) DeleteItWatcher(path string) {
	if _, ok := common.WatcherMap[path]; ok {
		err := notifyFile.watch.Remove(path)
		if err != nil {
			log.Println(err)
			return
		}
		delete(common.WatcherMap, path)
	}
}

// AppendChangedToOutputFile 追加变更历史
func AppendChangedToOutputFile(typeStr string, fileName string, isDir bool) {
	file, err := os.OpenFile(common.OutputFileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
	}(file)

	writer := csv.NewWriter(file)
	record := []string{typeStr, common.FileType[isDir], fileName, time.Now().String()}
	writeErr := writer.Write(record)
	if writeErr != nil {
		log.Println(writeErr)
		return
	}
	writer.Flush()
}

func LetItChanged(typeId int, fileName string) {
	common.Locker.Lock()
	if _, ok := common.ChangedMap[typeId]; !ok {
		common.ChangedMap[typeId] = list.New()
	}
	_, ok := common.WatcherMap[fileName]
	common.ChangedMap[typeId].PushBack(fileName)
	AppendChangedToOutputFile(common.PrefixMap[typeId], fileName, ok)
	common.Locker.Unlock()
}

// OutPutToFile 输出到指定文件中
func OutPutToFile() {
	f, err := os.Create(common.OutputFileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入utf8-bom
	writer := csv.NewWriter(f)

	for k, v := range common.Md5Map {
		isDir := false
		fileInfo, err := os.Stat(k)
		if err == nil && fileInfo.IsDir() {
			isDir = true
		}
		writer.Write([]string{common.FileType[isDir], k, v})
	}
	writer.Write([]string{""})
	writer.Write([]string{"文件变更历史："})
	writer.Flush()
}
