package config

import (
	"testing"
)

func TestGetConfig(t *testing.T) {
	ConfigFile = "../../../conf/config.yaml"

	InitConfig()

	setting := ServerSetting
	t.Log(setting.Sftp.Hostname)
}