package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
	"time"
)

type Config struct {
	Env        string     `yaml:"env" env-default:"local"`
	HTTPServer HTTPServer `yaml:"http_server" env-required:"true"`
	DBConfig   DBConfig   `yaml:"db_config" env-required:"true"`
}

type DBConfig struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-required:"true"`
	Username string `yaml:"user" env:"DB_USER" env-required:"true"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-required:"true"`
	Name     string `yaml:"name" env-required:"true"`
}

type HTTPServer struct {
	Address           string        `yaml:"address" env-default:"0.0.0.0:8080"`
	ReadTimeout       time.Duration `yaml:"read_timeout" env-default:"4s"`
	WriteTimeout      time.Duration `yaml:"write_timeout" env-default:"8s"`
	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout" env-default:"2s"`
	IdleTimeout       time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoadConfig() *Config {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("Variable CONFIG_PATH not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("Config file does not exist")
	}

	var config Config
	err = cleanenv.ReadConfig(configPath, &config)

	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	//todo CONFIGURE DATABASE

	return &config
}

func (config *Config) HTTPListen(e *echo.Echo) error {
	server := &http.Server{
		Addr:              config.HTTPServer.Address,
		ReadTimeout:       config.HTTPServer.ReadTimeout,
		WriteTimeout:      config.HTTPServer.WriteTimeout,
		ReadHeaderTimeout: config.HTTPServer.ReadHeaderTimeout,
		IdleTimeout:       config.HTTPServer.IdleTimeout,
		Handler:           e,
	}

	return e.StartServer(server)
}
