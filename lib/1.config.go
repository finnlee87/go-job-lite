package lib

import (
	"github.com/spf13/viper"
)

var Config *viper.Viper

func init() {
	// Default
	Config = viper.New()
	Config.SetConfigName("default")
	Config.AddConfigPath("config")
	Config.ReadInConfig()
	// Merge according to RUN_MODE
	Config.SetEnvPrefix("run")
	Config.AutomaticEnv()
	mode := Config.GetString("mode")
	if len(mode) > 0 {
		Config.SetConfigName(mode)
		Config.MergeInConfig()
	}
}
