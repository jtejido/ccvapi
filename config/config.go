package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

// This is the base Config type for the API. Extend as needed.
type Config struct {
	Title string
	Http  APIHTTP `toml:"http"`
}

type APIHTTP struct {
	Host          int    `toml:"host"`
	AccessLog     string `toml:"access_log_path"`
	ErrorLog      string `toml:"error_log_path"`
	CardTypesPath string `toml:"card_types_path"`
}

func LoadConfig(filename string) (*Config, error) {
	if filename == "" {
		return nil, fmt.Errorf("Unable to open file from an empty path")
	}

	var c Config

	if _, err := toml.DecodeFile(filename, &c); err != nil {
		return nil, err
	}

	return &c, nil
}
