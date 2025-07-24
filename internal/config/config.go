package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	HTTPServer `yaml:"httpServer"`
	DBPostgres `yaml:"dbPostgres"`
}

type HTTPServer struct {
	Addr         string        `yaml:"addr" env-default:":8080"`
	ReadTimeout  time.Duration `yaml:"readTimeout" env-default:"4s"`
	WriteTimeout time.Duration `yaml:"writeTimeout" env-default:"4s"`
	IdleTimeout  time.Duration `yaml:"idleTimeout" env-default:"60s"`
}

type DBPostgres struct {
	Host        string `yaml:"host" env-required:"true"`
	Port        int    `yaml:"port" env-required:"true"`
	User        string `yaml:"user" env-required:"true"`
	Password    string `yaml:"password" env-required:"true"`
	Db          string `yaml:"db" env-required:"true"`
	PoolMaxConn int    `yaml:"poolMaxConn" env-default:"4"`
}

func (d DBPostgres) Url() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?pool_max_conns=%d",
		url.QueryEscape(d.User),
		url.QueryEscape(d.Password),
		d.Host,
		d.Port,
		d.Db,
		d.PoolMaxConn,
	)
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
