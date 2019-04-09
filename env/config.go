package env

import (
	"github.com/spf13/viper"
	"log"
)

type Configuration struct {
	Cache Cache
	Server Server
}

type Server struct {
	Port int
}

type Cache struct {
	MaxSize string
	MaxItemSize string
	Duration string
}

func GetConfig() Configuration {
	var config Configuration
	viper.SetConfigName("conf")
	viper.AddConfigPath("../cassius")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	return config
}