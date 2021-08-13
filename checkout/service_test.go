package checkout

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gitlab.com/lucafmarques/hash-test/config"
	"gitlab.com/lucafmarques/hash-test/mocks"
)

func TestNewCheckoutService(t *testing.T) {
	config := config.ServiceConfig{
		Port:        ":TEST",
		Environment: config.DEVELOPMENT,
	}
	client := mocks.MockDiscountClient{}
	repo := mocks.MockRepository{}

	svc := NewCheckoutService(config, client, repo)

	assert.IsType(t, &Service{}, svc, "Failed asserting type of Service")
}

func TestNewCheckoutServiceDefaults(t *testing.T) {
	config := config.ServiceConfig{}
	client := mocks.MockDiscountClient{}
	repo := mocks.MockRepository{}

	svc := NewCheckoutService(config, client, repo)

	assert.IsType(t, &Service{}, svc, "Failed asserting type of Service")
}

func TestApplyMiddleware(t *testing.T) {
	server := echo.New()

	svc := Service{
		Server: *server,
	}

	buf := new(bytes.Buffer)

	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	}
	svc.Server.GET("/test", handler)

	middleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			buf.WriteString("called")
			return next(c)
		}
	}

	svc.ApplyMiddlewares(middleware)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	svc.Server.ServeHTTP(rec, req)

	assert.Equal(t, "called", buf.String(), "Failed asserting middleware was registered and called in request")
}

func TestRegisterRoutes(t *testing.T) {
	server := echo.New()

	svc := Service{
		Server: *server,
		Config: config.ServiceConfig{
			Environment: config.DEVELOPMENT,
		},
	}

	assert.Equal(t, len(svc.Server.Routes()), 0, "Failed asserting no routes have been registered yet")

	svc.RegisterRoutes()
	assert.Greater(t, len(svc.Server.Routes()), 0, "Failed asserting any routes have been registered yet")
}
