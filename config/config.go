package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var Config appConfig

type appConfig struct {
	CacheExpiryDuration  int `mapstructure:"cache_expiry_duration"`
	CacheCleanupInterval int `mapstructure:"cache_cleanup_interval"`
	ServerPort           int `mapstructure:"server_port"`
	APIDefaultTimeout    int `mapstructure:"api_default_timeout"`
}

func LoadConfig(configPaths ...string) error {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.SetEnvPrefix("")
	v.AutomaticEnv()

	for _, path := range configPaths {
		v.AddConfigPath(path)
	}

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to load config file %s", err)
	}

	// Reviewcomments --- dont leave commented code ---

	// Config.CacheExpiryDuration = v.Get("CacheExpiryDuration").(int)
	// Config.CacheCleanupInterval = v.Get("CacheCleanupInterval").(int)
	// Config.APIDefaultTimeout = v.Get("APIDefaultTimeout").(int)
	// Config.ServerPort = v.Get("ServerPort").(int)

	//v.SetDefault("ServerPort", Config.ServerPort)

	return v.Unmarshal(&Config)
}
