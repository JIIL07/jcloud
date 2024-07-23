package config

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/JIIL07/cloudFiles-manager/internal/storage"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string `yaml:"env" env-required:"true"`
	Server   ServerConfig
	Database storage.DBConfig
	URL      string `yaml:"url" env-required:"true"`
}

type Connection struct {
	IP net.IP
}

type ServerConfig struct {
	Address      string        `yaml:"address" env-default:"8080"`
	ReadTimeout  time.Duration `yaml:"readTimeout" env-default:"5s"`
	WriteTimeout time.Duration `yaml:"writeTimeout" env-default:"10s"`
	IdleTimeout  time.Duration `yaml:"idleTimeout" env-default:"15s"`
}

type User struct {
	Userame  string
	Password string
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
