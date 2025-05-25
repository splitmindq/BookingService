package main

import (
	"BookingService/internal/config"
	"BookingService/internal/controller"
	mw "BookingService/internal/middleware"
	"BookingService/internal/service"
	"BookingService/internal/storage/pgx"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envProd  = "prod"
	envDev   = "dev"
)

func main() {
	cfg := config.MustLoadConfig()
	logger := setupLogger(cfg.Env)

	storage, err := pgx.NewStorage(cfg)
	if err != nil {
		logger.Error("failed to init storage", "error", err)
		os.Exit(1)
	}
	defer storage.Close()

	userRepo := pgx.NewUserRepo(storage.GetPool(), logger)
	authService := service.NewAuthService(userRepo, cfg, logger)
	authController := controller.NewAuthController(authService, logger)
	jwtMiddleware := mw.NewJwtAuthMiddleware(userRepo, cfg.HTTPServer.JwtSecret, logger)

	e := echo.New()

	// CORS middleware

	// Логирование запросов

	api := e.Group("/api/v1")
	{
		api.POST("/sign_up", authController.SignUp())
		api.POST("/sign_in", authController.SignIn())

		// Защищенные роуты
		protected := api.Group("")
		protected.Use(jwtMiddleware.JwtAuth())
		{
			protected.GET("/profile", func(c echo.Context) error {
				userID, ok := c.Get("user_id").(int64)
				if !ok {
					logger.Warn("failed to get user_id from context")
					return echo.NewHTTPError(http.StatusUnauthorized, "invalid user context")
				}
				role, ok := c.Get("role").(string)
				if !ok {
					role = "unknown"
				}

				logger.Info("Profile accessed",
					slog.Int64("user_id", userID),
					slog.String("role", role))

				return c.JSON(http.StatusOK, map[string]interface{}{
					"user_id": userID,
					"role":    role,
					"message": "Protected content",
				})
			})
		}

		// Админские роуты
		admin := api.Group("/admin")
		admin.Use(jwtMiddleware.JwtAuth(), jwtMiddleware.AdminOnly())
		{
			admin.GET("/dashboard", func(c echo.Context) error {
				userID, _ := c.Get("user_id").(int64)
				logger.Info("Admin dashboard accessed", slog.Int64("user_id", userID))
				return c.JSON(http.StatusOK, map[string]string{
					"message": "Admin dashboard",
				})
			})
		}
	}

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
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case envDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return logger
}
