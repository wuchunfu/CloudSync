package config

import (
	"github.com/spf13/cobra"
	"github.com/wuchunfu/CloudSync/api"
	"github.com/wuchunfu/CloudSync/middleware/configx"
)

var StartCmd = &cobra.Command{
	Use:          "config",
	SilenceUsage: true,
	Short:        "Get Application config info",
	Example:      "CloudSync config -f conf/config.yaml",
	Run: func(cmd *cobra.Command, args []string) {
		api.Run()
	},
}

func init() {
	cobra.OnInitialize(configx.InitConfig)

	setting := configx.ServerSetting

	StartCmd.PersistentFlags().StringVarP(&configx.ConfigFile, "configFile", "f", "conf/config.yaml", "config file")
	StartCmd.PersistentFlags().StringVarP(&setting.Sftp.Hostname, "hostname", "H", "127.0.0.1", "hostname")
	StartCmd.PersistentFlags().IntVarP(&setting.Sftp.SshPort, "sshPort", "P", 22, "ssh port")
	StartCmd.PersistentFlags().StringVar(&setting.Sftp.Username, "username", "u", "username")
	StartCmd.PersistentFlags().StringVar(&setting.Sftp.Password, "password", "p", "password")
	// 必须配置项
	_ = StartCmd.MarkFlagRequired("configFile")

	// 使用viper可以绑定flag
	_ = configx.Vip.BindPFlag("sftp.hostname", StartCmd.PersistentFlags().Lookup("hostname"))
	_ = configx.Vip.BindPFlag("sftp.sshPort", StartCmd.PersistentFlags().Lookup("sshPort"))
	_ = configx.Vip.BindPFlag("sftp.username", StartCmd.PersistentFlags().Lookup("username"))
	_ = configx.Vip.BindPFlag("sftp.password", StartCmd.PersistentFlags().Lookup("password"))

	// 设置默认值
	configx.Vip.SetDefault("sftp.username", "root")
	configx.Vip.SetDefault("sftp.sshPort", "22")
}
