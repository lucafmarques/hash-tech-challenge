package checkout

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

var (
	ErrInvalidProductId       = echo.NewHTTPError(http.StatusBadRequest, "Product ID should be a number.")
	ErrInvalidPayload         = echo.NewHTTPError(http.StatusBadRequest, "Request body doesn't comply with API spec.")
	ErrFailedFetchingProducts = echo.NewHTTPError(http.StatusInternalServerError, "Couldn't fetch all products.")
)

// @summary Calculate Checkout
// @description Calculate Checkout for cart of products
// @accept json
// @produce json
// @param cart body CheckoutRequest true "Cart of products"
// @success 200 {object} CheckoutResponse
// @failure 400 {object} HTTPError
// @router /checkout [post]
func (svc *Service) PostCheckout(c echo.Context) error {
	ctx := c.Request().Context()

	data := CheckoutRequest{}
	if err := c.Bind(&data); err != nil {
		log.Warnf("Failed to unmarshal reques data into %T: %v", data, err)
		return ErrInvalidPayload
	}

	discounts, totalAmount, totalDiscount := svc.Core.CalculateCheckout(ctx, data.Products)

	gift, ok := svc.Core.BlackFridayGift()
	if ok {
		discounts = append(discounts, *gift)
	}

	response := CheckoutResponse{
		Products:                discounts,
		TotalAmount:             totalAmount,
		TotalDiscount:           totalDiscount,
		TotalAmountWithDiscount: totalAmount - totalDiscount,
	}

	return c.JSON(http.StatusOK, response)
}
