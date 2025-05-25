package middleware

import (
	"BookingService/internal/lib/jwt"
	"BookingService/internal/repository"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"strings"
)

type JwtAuthMiddleware struct {
	userRepo repository.UserRepository
	secret   string
	logger   *slog.Logger
}

func NewJwtAuthMiddleware(userRepo repository.UserRepository, secret string, logger *slog.Logger) *JwtAuthMiddleware {
	return &JwtAuthMiddleware{userRepo, secret, logger}
}

func (m *JwtAuthMiddleware) JwtAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			// Проверяем заголовок Authorization
			authHeader := ctx.Request().Header.Get("Authorization")
			tokenString := ""
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					tokenString = parts[1]
				}
			}

			// Если заголовок пустой, проверяем куки
			if tokenString == "" {
				cookie, err := ctx.Cookie("token")
				if err != nil {
					m.logger.Warn("missing token")
					return echo.NewHTTPError(http.StatusUnauthorized, "missing token")
				}
				tokenString = cookie.Value
			}

			if tokenString == "" {
				m.logger.Warn("missing token")
				return echo.NewHTTPError(http.StatusUnauthorized, "missing token")
			}

			userID, role, err := jwt.ParseTokenAndGetUID(tokenString, m.secret)
			if err != nil {
				m.logger.Warn("invalid token", "error", err)
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
			}

			// Сохраняем user_id и role в контексте
			ctx.Set("user_id", userID)
			ctx.Set("role", role)

			m.logger.Debug("Token validated",
				slog.Int64("user_id", userID),
				slog.String("role", role))

			return next(ctx)
		}
	}
}

func (m *JwtAuthMiddleware) AdminOnly() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			role, ok := ctx.Get("role").(string)
			if !ok || role != "admin" {
				m.logger.Warn("admin access required", "role", role)
				return echo.NewHTTPError(http.StatusForbidden, "admin access required")
			}
			return next(ctx)
		}
	}
}
