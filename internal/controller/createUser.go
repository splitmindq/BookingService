package controller

import (
	"BookingService/internal/entity"
	"BookingService/internal/util"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
)

func NewUser(log *slog.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		const op = "control/NewUser"
		log = log.With(
			slog.String("op", op),
			slog.String("method", c.Request().Method),
			slog.String("path", c.Path()),
		)

		var user entity.User
		if err := util.DecodeAndValidate(c, &user); err != nil {
			return err
		}

		log.Info("user created successfully",
			slog.String("name", user.Name),
			slog.String("surname", user.Surname),
			slog.String("phone", user.Phone))

		return c.JSON(http.StatusCreated, map[string]string{
			"status": "user created",
		})

	}
}
