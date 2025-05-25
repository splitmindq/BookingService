package util

import (
	valid "BookingService/internal/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
)

var validate = valid.Validate

func DecodeAndValidate(c echo.Context, v interface{}) error {
	if err := c.Bind(v); err != nil {
		log.Error("failed to bind data", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request payload")
	}

	if err := validate.Struct(v); err != nil {
		log.Error("failed to validate data", err)
		return echo.NewHTTPError(http.StatusBadRequest, "validation error")
	}

	return nil
}
