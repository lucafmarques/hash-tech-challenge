package auth

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	api_key := "TEST_API_KEY"
	APIKeys = []string{api_key}
	os.Setenv(KeyEnv, api_key)

	middlewareFunc := AuthMiddleware()

	handlerCalled := false
	handler := func(c echo.Context) error {
		handlerCalled = true
		return c.String(http.StatusOK, "test")
	}
	middlewareChain := middlewareFunc(handler)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", api_key))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := middlewareChain(c)

	assert.NoError(t, err, "Failed asserting no err in AuthMiddleware chain")
	assert.True(t, handlerCalled)
}

func TestAuthMiddlewareInvalid(t *testing.T) {
	middlewareFunc := AuthMiddleware()

	handlerCalled := false
	handler := func(c echo.Context) error {
		handlerCalled = true
		return c.String(http.StatusOK, "test")
	}
	middlewareChain := middlewareFunc(handler)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", "INVALID_TEST_API_KEY"))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := middlewareChain(c)

	assert.Error(t, err, "Failed asserting err in AuthMiddleware chain")
	assert.False(t, handlerCalled)
}
