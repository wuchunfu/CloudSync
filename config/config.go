package config

import (
	"encoding/json"
	"os"
)

// Global 定义全局常量 用户根据 conf 里的文件 conf.json 来配置
type Global struct {
	Name    string
	Version string
	Host    string
	Sftp    GlobalSftpMap
	Sync    []GlobalSyncMap
	LogPath string
}

// GlobalSftpMap sftp当前服务器主机
type GlobalSftpMap struct {
	Hostname string
	Username string
	Password string
	SSHPort  int
}

// GlobalSyncMap 同步路径参数
type GlobalSyncMap struct {
	Name       string
	SourcePath string
	TargetPath string
}

// GlobalObject 全局配置
var GlobalObject *Global

//Reload 读取用户的配置文件
func (global *Global) Reload() {
	data, err := os.ReadFile("/Users/wuchunfu/workspace/learning/codespace/CloudSync/conf/conf.json")
	if err != nil {
		panic(err)
	}
	// 将json数据解析到struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

// 提供init方法, 默认加载
func init() {
	// 初始化GlobalObject变量,设置一些默认值
	GlobalObject = &Global{
		Name:    "ServerApp",
		Version: "V1.0",
		Host:    "0.0.0.0",
	}
	// 从配置文件中加载一些用户配置的参数
	GlobalObject.Reload()
}
