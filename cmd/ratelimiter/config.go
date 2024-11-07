package main

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Redis RedisConf
	GRPC  GRPCConf
	Rates RatesConf
	DB    DBConf
}

type RedisConf struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type GRPCConf struct {
	Port string
}

type RatesConf struct {
	Login    int
	Password int
	IP       int
}

type DBConf struct {
	User     string
	Password string
	Name     string
	Host     string
	Port     string
}

func LoadConfig(path string) (Config, error) {
	viper.SetConfigFile(path)

	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, fmt.Errorf("failed reading config: %w", err)
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var config Config
	viper.Unmarshal(&config)
	if err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return config, nil
}
