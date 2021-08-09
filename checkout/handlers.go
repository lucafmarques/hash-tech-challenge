package checkout

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

var (
	ErrInvalidProductId       = echo.NewHTTPError(http.StatusBadRequest, "Product ID should be a number.")
	ErrInvalidPayload         = echo.NewHTTPError(http.StatusBadRequest, "Request body doesn't comply with API spec.")
	ErrFailedFetchingProducts = echo.NewHTTPError(http.StatusInternalServerError, "Couldn't fetch all products.")
)

func (svc *Service) GetAllProducts(c echo.Context) error {
	products, err := svc.Core.AllProducts()
	if err != nil {
		log.Warnf("Failed fetching all products from repository: %v", err)
		return ErrFailedFetchingProducts
	}
	return c.JSON(http.StatusOK, products)
}

func (svc *Service) GetProductDiscount(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.Atoi((c.Param("id")))
	if err != nil {
		log.Warnf("Failed conversion of Product ID to int: %v", err)
		return ErrInvalidProductId
	}

	percentage := svc.Core.CalculateDiscountPercentage(ctx, id)

	resp := DiscountResponse{
		Percentage: percentage,
	}

	return c.JSON(http.StatusOK, resp)
}

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
