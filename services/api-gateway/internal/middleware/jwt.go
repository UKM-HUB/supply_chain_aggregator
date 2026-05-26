package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func JWT(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authorization := c.Request().Header.Get("Authorization")
			if authorization == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "authorization header is required",
				})
			}

			if !strings.HasPrefix(authorization, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "authorization header must use Bearer token",
				})
			}

			token := strings.TrimSpace(strings.TrimPrefix(authorization, "Bearer "))
			if token == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "token is required",
				})
			}

			_ = secret
			c.Set("access_token", token)

			return next(c)
		}
	}
}
