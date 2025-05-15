package main

import (
	"BookingService/internal/config"
	"BookingService/internal/controller"
	"BookingService/internal/service"
	"BookingService/internal/storage/pgx"
	"fmt"
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
	fmt.Println(cfg)
	logger := setupLogger(cfg.Env)

	storage, err := pgx.NewStorage(cfg)
	if err != nil {
		logger.Error(err.Error())
	}
	defer storage.Close()

	userRepo := pgx.NewUserRepo(storage.GetPool(), logger)

	authService := service.NewAuthService(userRepo, cfg, logger)

	authController := controller.NewAuthController(authService)

	e := echo.New()

	//тестовая ручка
	e.POST("/users", controller.NewUser(logger))

	//тест ручка без middleware
	e.POST("/signup", authController.SignUp())

	err = cfg.HTTPListen(e)
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
