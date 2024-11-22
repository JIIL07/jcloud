package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string `yaml:"env" env-required:"true"`
	Server   ServerConfig
	Database DBConfig
	Static   StaticPath
	URL      string `yaml:"url" env-required:"true"`
}

type ServerConfig struct {
	Address      string        `yaml:"address" env-default:"8080"`
	ReadTimeout  time.Duration `yaml:"readTimeout" env-default:"10s"`
	WriteTimeout time.Duration `yaml:"writeTimeout" env-default:"20s"`
	IdleTimeout  time.Duration `yaml:"idleTimeout" env-default:"60s"`
}

type DBConfig struct {
	DriverName     string `yaml:"driverName" env-required:"true"`
	DataSourceName string `yaml:"dataSourceName" env-required:"true"`
}

type StaticPath struct {
	Path string `yaml:"path" env-required:"true"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
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
