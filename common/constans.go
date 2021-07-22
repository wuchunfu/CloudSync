package common

import (
	"container/list"
	"sync"
)

var (
	Locker sync.Mutex

	Md5Map map[string]string

	ChangedMap map[int]*list.List

	WatcherMap map[string]bool

	PrefixMap = map[int]string{
		1: "新建",
		2: "修改",
		3: "删除",
		4: "重命名",
		5: "修改权限",
	}

	FileType = map[bool]string{
		true:  "文件夹",
		false: "文件",
	}

	OutputFileName = "filesName.csv"
)
