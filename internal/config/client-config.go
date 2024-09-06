package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type ClientConfig struct {
	Client        Client        `yaml:"client"`
	Services      Services      `yaml:"services"`
	RetryStrategy RetryStrategy `yaml:"retry_strategy"`
	Metrics       Metrics       `yaml:"metrics"`
}

type Client struct {
	ID          string    `yaml:"id"`
	Name        string    `yaml:"name"`
	Environment string    `yaml:"environment"`
	API         APIConfig `yaml:"api"`
	Logging     Logging   `yaml:"logging"`
}

type APIConfig struct {
	BaseURL   string `yaml:"base_url"`
	AuthToken string `yaml:"auth_token"`
	Timeout   int    `yaml:"timeout"`
	Retries   int    `yaml:"retries"`
	SSLVerify bool   `yaml:"ssl_verify"`
}

type Logging struct {
	Level    string `yaml:"level"`
	LogPath  string `yaml:"log_path"`
	MaxSize  int    `yaml:"max_size"`
	MaxFiles int    `yaml:"max_files"`
}

type Services struct {
	Database      DatabaseConfig `yaml:"database"`
	Cache         CacheConfig    `yaml:"cache"`
	MessageBroker MessageBroker  `yaml:"message_broker"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SSLMode  string `yaml:"sslmode"`
}

type CacheConfig struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	UseSSL bool   `yaml:"use_ssl"`
}

type MessageBroker struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	QueueName string `yaml:"queue_name"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	UseSSL    bool   `yaml:"use_ssl"`
}

type RetryStrategy struct {
	MaxAttempts       int     `yaml:"max_attempts"`
	DelaySeconds      int     `yaml:"delay_seconds"`
	BackoffMultiplier float64 `yaml:"backoff_multiplier"`
}

type Metrics struct {
	Enabled         bool   `yaml:"enabled"`
	Provider        string `yaml:"provider"`
	Endpoint        string `yaml:"endpoint"`
	IntervalSeconds int    `yaml:"interval_seconds"`
}

func MustLoadClient() *ClientConfig {
	configPath := os.Getenv("CLIENT_CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %v", configPath)
	}

	var cfg ClientConfig

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %v", err)
	}

	return &cfg
}
