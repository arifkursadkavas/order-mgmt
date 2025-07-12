package config

import (
	"fmt"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

var Config appConfig

type appConfig struct {
	DBClient *mongo.Client

	DBConnectionString string `mapstructure:"dbConnectionString"`

	DBDefaultTimeout int

	UserSessionTimeout int64

	ServerPort int `mapstructure:"server_port"`

	ApiKey string `mapstructure:"api_key"`
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
	Config.DBDefaultTimeout = v.Get("DBDefaultTimeout").(int)

	Config.ApiKey = v.Get("API_KEY").(string)
	Config.ServerPort = v.Get("ServerPort").(int)

	v.SetDefault("ServerPort", Config.ServerPort)

	return v.Unmarshal(&Config)
}
