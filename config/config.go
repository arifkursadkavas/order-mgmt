package config

import (
	"fmt"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

var Config appConfig

type appConfig struct {
	DBClient           *mongo.Client
	DBConnectionString string `mapstructure:"dbConnectionString"`
	ServerPort         int    `mapstructure:"server_port"`
	APIDefaultTimeout  int    `mapstructure:"api_key"`
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

	Config.DBConnectionString = v.Get("DBConnectionString").(string)

	Config.APIDefaultTimeout = v.Get("APIDefaultTimeout").(int)
	Config.ServerPort = v.Get("ServerPort").(int)

	v.SetDefault("ServerPort", Config.ServerPort)

	return v.Unmarshal(&Config)
}
