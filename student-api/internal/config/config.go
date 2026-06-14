package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HttpServer struct {
	Addr string
}

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true" env-default:"production"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HttpServer  `yaml:"http-server"`
}

func MustLoad() *Config {
	var ConfigPath string

	ConfigPath = os.Getenv("CONFIG_PATH")

	if ConfigPath == "" {
		flags := flag.String("config", "", "Path to the config file")
		flag.Parse()

		ConfigPath = *flags
	}

	if ConfigPath == "" {
		log.Fatal("Config is not set")
	}

	if _, err := os.Stat(ConfigPath); os.IsNotExist(err) {
		log.Fatalf("Config file is not exist %s", ConfigPath)
	}

	var cfg Config

	err := cleanenv.ReadConfig(ConfigPath, &cfg)
	if err != nil {
		log.Fatalf("cannot read config file: %s", err.Error())
	}

	return &cfg

}
