package middleware

import (
	"github.com/LightningFootball/backend/app/response"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func ValidateParams(intParams map[string]string) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			for p, v := range intParams {
				param := c.Param(p)
				if param != "" {
					if _, err := strconv.Atoi(param); err != nil {
						return c.JSON(http.StatusNotFound, response.ErrorResp(v, nil))
					}
				}
			}
			return next(c)
		}
	}
}
