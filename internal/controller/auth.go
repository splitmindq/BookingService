package controller

import (
	"BookingService/internal/entity"
	"BookingService/internal/service"
	"errors"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"time"
)

type AuthController struct {
	service *service.AuthService
	logger  *slog.Logger
}

func NewAuthController(service *service.AuthService, logger *slog.Logger) *AuthController {
	return &AuthController{
		service: service,
		logger:  logger,
	}
}

func (c *AuthController) SignUp() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var input entity.SignUpInput
		if err := ctx.Bind(&input); err != nil {
			c.logger.Warn("failed to bind signup input", "error", err)
			return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
		}

		userID, err := c.service.SignUp(ctx.Request().Context(), input)
		if err != nil {
			c.logger.Error("signup failed", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "registration failed")
		}

		return ctx.JSON(http.StatusCreated, map[string]int64{
			"user_id": userID,
		})
	}
}

func (c *AuthController) SignIn() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var input entity.SignInInput
		if err := ctx.Bind(&input); err != nil {
			c.logger.Warn("failed to bind signin input", "error", err)
			return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
		}

		token, err := c.service.SignIn(ctx.Request().Context(), input)
		if err != nil {
			errCreds := errors.New("invalid credentials")
			if errors.Is(err, errCreds) {
				c.logger.Warn("invalid login attempt", "email", input.Email)
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
			}
			c.logger.Error("signin failed", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}

		cookie := new(http.Cookie)
		cookie.Name = "token"
		cookie.Value = token
		cookie.Expires = time.Now().Add(c.service.Cfg.HTTPServer.JwtExpire)
		cookie.Path = "/"
		cookie.HttpOnly = true
		cookie.Secure = true
		cookie.SameSite = http.SameSiteStrictMode

		ctx.SetCookie(cookie)

		return ctx.JSON(http.StatusOK, map[string]string{
			"token": token,
		})
	}
}
