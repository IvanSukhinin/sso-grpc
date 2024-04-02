package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env            string     `yaml:"env" env-required:"true"`
	GRPC           GRPCConfig `yaml:"grpc"`
	MigrationsPath string
	TokenTTL       time.Duration `yaml:"token_ttl" env-default:"1h"`
	Db             `yaml:"db"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type Db struct {
	Host     string `yaml:"host" env-required:"true" env:"POSTGRES_HOST"`
	Port     string `yaml:"port" env-required:"true" env:"POSTGRES_PORT"`
	User     string `yaml:"user" env-required:"true" env:"POSTGRES_USER"`
	Password string `yaml:"password" env-required:"true" env:"POSTGRES_PASSWORD"`
	Db       string `yaml:"db" env-required:"true" env:"POSTGRES_DB"`
}

func MustLoad() *Config {
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
