package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App struct {
		Env      string `yaml:"env"`
		LogLevel string `yaml:"logLevel"`
	} `yaml:"app"`
	HttpServer struct {
		Host            string `yaml:"host"`
		Port            string `yaml:"port"`
		ReadTimeout     int16  `yaml:"readTimeout"`
		WriteTimeout    int16  `yaml:"writeTimeout"`
		ShutdownTimeout int16  `yaml:"shutdownTimeout"`
	} `yaml:"httpServer"`
	Repository struct {
		Postgres struct {
			Port     uint16 `yaml:"port"`
			Host     string `yaml:"host"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
			Database string `yaml:"database"`
		} `yaml:"postgres"`
	} `yaml:"repository"`
}

var (
	cfg  Config
	once sync.Once
)

func Read(configPath string) {
	once.Do(func() {
		err := cleanenv.ReadConfig(configPath, &cfg)
		if err != nil {
			log.Fatalf("Failed to read config: %v", err)
		}
	})
}

func Get() *Config {
	return &cfg
}
