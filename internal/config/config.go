package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string        `yaml:"env" env-default:"local" env-required:"true"`
	HTTPServer HTTPServer    `yaml:"http_server"`
	PosgreSQL  PosgreSQL     `yaml:"postgres"`
	Clients    ClientsConfig `yaml:"clients"`
	AppSecret  string        `yaml:"app_secret" env-required:"true" env"APP_SECRET"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	User        string        `yaml:"user" env-required:"true"`
	Password    string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

type PosgreSQL struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type Client struct {
	Addres       string        `yaml:"address"`
	Timeout      time.Duration `yaml:"timeout"`
	RetriesCount int           `yaml:"retries_count"`
	// Insecure     bool          `yaml:"insecure"`
}

type ClientsConfig struct {
	SSO Client `yaml:"sso"`
}

// TODO: postgress part for config

func Load() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PARH is not set")
	}

	if _, err := os.Stat(configPath); os.IsExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	return &cfg
}
