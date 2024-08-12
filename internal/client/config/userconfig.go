package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

var u string

func ReadUserConfig() error {
	f, err := os.OpenFile(os.Getenv("USER_CONFIG_PATH"), os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	return yaml.NewDecoder(f).Decode(&u)
}
func WriteUserConfig() error {
	f, err := os.OpenFile(os.Getenv("USER_CONFIG_PATH"), os.O_RDONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	return yaml.NewEncoder(f).Encode(u)
}
