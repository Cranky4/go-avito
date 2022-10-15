package main

import (
	"log"

	internalsender "github.com/Cranky4/go-avito/hw12_13_14_15_calendar/internal/sender"
	"github.com/spf13/viper"
)

func NewConfig(path string) internalsender.Config {
	viper.SetConfigFile(path)
	var c internalsender.Config

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	if err := viper.Unmarshal(&c); err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}

	return c
}
