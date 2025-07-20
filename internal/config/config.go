package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	HTTPServer `yaml:"httpServer"`
	DBServer   `yaml:"dbServer"`
}

type HTTPServer struct {
	Addr         string        `yaml:"addr" env-default:":8080"`
	ReadTimeout  time.Duration `yaml:"readTimeout" env-default:"4s"`
	WriteTimeout time.Duration `yaml:"writeTimeout" env-default:"4s"`
	IdleTimeout  time.Duration `yaml:"idleTimeout" env-default:"60s"`
}

type DBServer struct {
	Dns string `yaml:"dns" env-required:"true"`
}

func Load() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
