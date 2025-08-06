package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string                    `yaml:"env" env-default:"local"`
	HTTPServer HTTPServer                `yaml:"httpServer"`
	Databases  map[string]ConfigDatabase `yaml:"databases"`
	Logger     Logger                    `yaml:"logger"`
}

type HTTPServer struct {
	Addr         string        `yaml:"addr" env-default:":8080"`
	ReadTimeout  time.Duration `yaml:"readTimeout" env-default:"4s"`
	WriteTimeout time.Duration `yaml:"writeTimeout" env-default:"4s"`
	IdleTimeout  time.Duration `yaml:"idleTimeout" env-default:"60s"`
}

type ConfigDatabase struct {
	URL        string `yaml:"url" env-required:"true"`
	MigrateUrl string `yaml:"migrate"`
}

type Logger struct {
	FilePath string `yaml:"filePath"`
	Level    string `yaml:"level" env-default:"Info"`
}

func Load() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		return nil, fmt.Errorf("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("cannot read config: %s", err)
	}

	return &cfg, nil
}
