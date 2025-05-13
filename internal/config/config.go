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
	Env        string        `yaml:"env" env-default:"local"`
	HTTPServer HTTPServer    `yaml:"http_server" env-required:"true"`
	DBConfig   DBConfig      `yaml:"db_config" env-required:"true"`
	JwtSecret  []byte        `env-required:"true"`
	JwtExpire  time.Duration `yaml:"jwt_expire" env-default:"1h"`
}

type DBConfig struct {
	Host                  string        `yaml:"host" env:"DB_HOST" env-required:"true"`
	Port                  int           `yaml:"port" env:"DB_PORT" env-required:"true"`
	Username              string        `yaml:"user" env:"DB_USER" env-required:"true"`
	Password              string        `yaml:"password" env:"DB_PASSWORD" env-required:"true"`
	Name                  string        `yaml:"name" env:"DB_NAME" env-required:"true"`
	MaxConnections        int32         `yaml:"max_connections" env-required:"true"`
	MinConnections        int32         `yaml:"min_connections" env-required:"true"`
	MaxConnectionLifetime time.Duration `yaml:"max_connection_life" env-required:"true"`
	MaxConnectionIdleTime time.Duration `yaml:"max_connection_idle_time" env-required:"true"`
}

type HTTPServer struct {
	Address           string        `yaml:"address" env-default:"0.0.0.0:8080"`
	ReadTimeout       time.Duration `yaml:"read_timeout" env-default:"4s"`
	WriteTimeout      time.Duration `yaml:"write_timeout" env-default:"8s"`
	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout" env-default:"2s"`
	IdleTimeout       time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoadConfig() *Config {

	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Fatalf("Error loading .env file: %v", err)
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("Environment variable CONFIG_PATH not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist at path: %s", configPath)
	}

	var config Config
	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	config.JwtSecret = []byte(os.Getenv("JWT_SECRET"))

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
