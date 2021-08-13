package checkout

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gitlab.com/lucafmarques/hash-test/config"
	"gitlab.com/lucafmarques/hash-test/discount"
	"gitlab.com/lucafmarques/hash-test/mocks"
	"gitlab.com/lucafmarques/hash-test/repository"
)

var (
	checkoutReq = `{"products": [{"id": 0, "quantity": 2}]}`
)

func TestCheckoutHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/checkout", strings.NewReader(checkoutReq))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	product := []*repository.Product{
		{
			ID:   0,
			Gift: true,
		},
	}
	client := mocks.MockDiscountClient{
		Resp: &discount.GetDiscountResponse{
			Percentage: 0.03,
		},
	}
	repo := mocks.MockRepository{
		Products: product,
		Product:  product[0],
	}
	svc := Service{
		Core: NewCore(config.CoreConfig{
			BlackFridayDate: time.Now().Format("01/02"),
		}, client, repo),
	}

	if assert.NoError(t, svc.PostCheckout(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestCheckoutHandlerInvalidRequest(t *testing.T) {
	checkoutReq := `{"products": [{"id": "test", "quantity": 2}]}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/checkout", strings.NewReader(checkoutReq))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	product := []*repository.Product{
		{
			ID: 0,
		},
	}
	client := mocks.MockDiscountClient{
		Resp: &discount.GetDiscountResponse{
			Percentage: 0.03,
		},
	}
	repo := mocks.MockRepository{
		Products: product,
		Product:  product[0],
	}
	svc := Service{
		Core: NewCore(config.CoreConfig{}, client, repo),
	}

	err := svc.PostCheckout(c)
	if assert.Error(t, err, "Failed asserting err") {
		httpErr, _ := err.(*echo.HTTPError)
		expected := http.StatusBadRequest
		assert.Equal(t, expected, httpErr.Code, fmt.Sprintf("Failed asserting %v status code", expected))
	}
}
