package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"regexp"
)

var Statuses = []string{"Created", "Has data in", "Renamed", "Deleted"}

var IsValidTableName = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_.]*$`).MatchString

type Config struct {
	Env      string `yaml:"env" env-required:"true"`
	Database DBConfig
}

type DBConfig struct {
	DriverName     string `yaml:"driverName" env-required:"true"`
	DataSourceName string `yaml:"dataSourceName" env-required:"true"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CLIENT_CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %v", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %v", err)
	}

	return &cfg
}
