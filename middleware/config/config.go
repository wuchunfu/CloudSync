package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wuchunfu/CloudSync/utils/fileUtils"
	"os"
)

// Sftp sftp server host
type Sftp struct {
	Hostname string `yaml:"hostname"`
	SshPort  int    `yaml:"sshPort"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// Sync sync path parameters
type Sync struct {
	Name       string `yaml:"name"`
	SourcePath string `yaml:"sourcePath"`
	TargetPath string `yaml:"targetPath"`
}

// YamlSetting global constants are defined and configured by the user according to the file conf.yaml in conf
type YamlSetting struct {
	Sftp        Sftp     `yaml:"sftp"`
	Sync        []Sync   `yaml:"sync"`
	IgnoreFiles []string `yaml:"ignoreFiles"`
}

var (
	Vip        = viper.New()
	ConfigFile = ""
	// ServerSetting global config
	ServerSetting = new(YamlSetting)
)

// InitConfig reads in config file and ENV variables if set.
func InitConfig() {
	if ConfigFile != "" {
		if !fileUtils.FilePathExists(ConfigFile) {
			logger.Errorf("No such file or directory: %s", ConfigFile)
			os.Exit(1)
		} else {
			// Use config file from the flag.
			Vip.SetConfigFile(ConfigFile)
			Vip.SetConfigType("yaml")
		}
	} else {
		logger.Errorf("Could not find config file: %s", ConfigFile)
		os.Exit(1)
	}
	// If a config file is found, read it in.
	err := Vip.ReadInConfig()
	if err != nil {
		logger.Errorf("Failed to get config file: %s", ConfigFile)
	}
	Vip.WatchConfig()
	Vip.OnConfigChange(func(e fsnotify.Event) {
		logger.Infof("Config file changed: %s\n", e.Name)
		fmt.Printf("Config file changed: %s\n", e.Name)
		ServerSetting = GetConfig(Vip)
	})
	Vip.AllSettings()
	ServerSetting = GetConfig(Vip)
}

// GetConfig 解析配置文件，反序列化
func GetConfig(vip *viper.Viper) *YamlSetting {
	setting := new(YamlSetting)
	// 解析配置文件，反序列化
	if err := vip.Unmarshal(setting); err != nil {
		logger.Errorf("Unmarshal yaml faild: %s", err)
		os.Exit(1)
	}
	return setting
}
