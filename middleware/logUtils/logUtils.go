package logUtils

import (
	"github.com/wuchunfu/CloudSync/config"
	"log"
	"os"
)

var Logger *log.Logger

func init() {
	file, err := os.Create(config.GlobalObject.LogPath)
	if err != nil {
		log.Fatalln("fail to create test.log file: ", err)
	}
	Logger = log.New(file, "", log.Llongfile)
	Logger.SetFlags(log.LstdFlags)
}
