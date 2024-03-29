package main

import (
	"log"

	sqlstorage "github.com/Cranky4/go-avito/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/spf13/viper"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger   LoggerConf
	Storage  StorageConf
	Database sqlstorage.DatabaseConf
	HTTP     HTTPConf
	GRPC     GrpcConf
}

type LoggerConf struct {
	Level string
}

type HTTPConf struct {
	Addr, RequestLogFile string
}

type StorageDriver string

type GrpcConf struct {
	Addr, RequestLogFile string
}

const (
	_          StorageDriver = "memory"
	SQLStorage StorageDriver = "sql"
)

type StorageConf struct {
	Driver StorageDriver
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
