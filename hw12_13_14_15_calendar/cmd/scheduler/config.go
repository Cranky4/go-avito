package main

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Logger   LoggerConf
	Database DatabaseConf
	Broker   BrokerConf
}

type LoggerConf struct {
	Level string
}

type DatabaseConf struct {
	Dsn string
}

type BrokerAdapter string

const (
	KafkaAdapter BrokerAdapter = "kafka"
)

type BrokerConf struct {
	Adapter BrokerAdapter
	Address string
}

func NewConfig(path string) Config {
	viper.SetConfigFile(path)
	var c Config

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	if err := viper.Unmarshal(&c); err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}

	return c
}
