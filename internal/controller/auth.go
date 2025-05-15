package controller

import (
	"BookingService/internal/entity"
	"BookingService/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

type AuthController struct {
	service *service.AuthService
	//logger
}

func NewAuthController(service *service.AuthService) *AuthController {
	return &AuthController{
		service: service,
	}
}
func (controller *AuthController) SignUp() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input entity.SignUpInput
		if err := c.Bind(&input); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		userID, err := controller.service.SignUp(c.Request().Context(), input)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusCreated, map[string]int64{
			"user_id": userID,
		})
	}
}
