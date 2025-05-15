package middleware

import (
	"BookingService/internal/lib/jwt"
	"context"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"strings"
)

type UserRepository interface {
	IsAdmin(ctx context.Context, userId int64) (bool, error)
}

type JwtAuthMiddleware struct {
	userRepo  UserRepository
	jwtSecret string
	log       *slog.Logger
}

func NewJwtAuthMiddleware(userRepo UserRepository, jwtSecret string, log *slog.Logger) *JwtAuthMiddleware {
	return &JwtAuthMiddleware{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		log:       log,
	}
}

func (jm *JwtAuthMiddleware) JwtAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing Authorization header")
			}
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing Authorization header")
			}
			uid, err := jwt.ParseTokenAndGetUID(tokenString, jm.jwtSecret)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
			}
			c.Set("uid", uid)
			return next(c)
		}
	}
}

func (jm *JwtAuthMiddleware) AdminOnly() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			uid, ok := c.Get("uid").(int64)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Not authorized")
			}
			IsAdmin, err := jm.userRepo.IsAdmin(context.Background(), uid)
			if err != nil {
				jm.log.Error("failed to check admin status", slog.Int64("uid", uid),
					slog.String("error", err.Error()))
				return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
			}
			if !IsAdmin {
				return echo.NewHTTPError(http.StatusUnauthorized, "Forbidden")
			}
			return next(c)
		}
	}
}
