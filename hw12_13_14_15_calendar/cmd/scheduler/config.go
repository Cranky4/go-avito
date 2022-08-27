package main

import (
	"log"

	internalscheduler "github.com/Cranky4/go-avito/hw12_13_14_15_calendar/internal/scheduler"
	"github.com/spf13/viper"
)

func NewConfig(path string) internalscheduler.Config {
	viper.SetConfigFile(path)
	var c internalscheduler.Config

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	if err := viper.Unmarshal(&c); err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}

	return c
}
