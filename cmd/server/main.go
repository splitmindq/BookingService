package main

import (
	"BookingService/internal/config"
	"BookingService/internal/controller"
	"github.com/labstack/echo/v4"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envProd  = "prod"
	envDev   = "dev"
)

func main() {

	//todo init logger
	//todo init storage
	//todo create routes

	cfg := config.MustLoadConfig()

	logger := setupLogger(cfg.Env)

	e := echo.New()

	e.POST("/users", controller.NewUser(logger))
	
	err := cfg.HTTPListen(e)
	if err != nil {
		e.Logger.Fatal(err)
	}

}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return logger
}
