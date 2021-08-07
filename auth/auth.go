package auth

import (
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var APIKeys []string

const KeyEnv = "API_KEYS"

func init() {
	APIKeys = strings.Split(os.Getenv(KeyEnv), ",")
}

func AuthMiddleware() echo.MiddlewareFunc {
	return middleware.KeyAuth(func(s string, c echo.Context) (bool, error) {
		for _, key := range APIKeys {
			if key == s {
				return true, nil
			}
		}
		return false, echo.NewHTTPError(http.StatusUnauthorized, "Invalid APIKey, contact service provider if issue persists.")
	})
}
