package checkout

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"gitlab.com/lucafmarques/hash-test/discount"
)

var (
	ErrInvalidProductId       = echo.NewHTTPError(http.StatusBadRequest, "Product ID should be a number.")
	ErrInvalidPayload         = echo.NewHTTPError(http.StatusBadRequest, "Request body doesn't comply with API spec.")
	ErrFailedFetchingProducts = echo.NewHTTPError(http.StatusInternalServerError, "Couldn't fetch all products.")
)

func (svc *Service) GetAllProducts(c echo.Context) error {
	products, err := svc.Repository.GetAllProducts()
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

	resp := DiscountResponse{
		Percentage: 0,
	}

	req := &discount.GetDiscountRequest{
		ProductID: int32(id),
	}

	discountResp, err := svc.DiscountClient.GetDiscount(ctx, req)
	if err != nil {
		log.Warnf("Failed requesting discount from external service: %v", err)
	} else {
		resp.Percentage = discountResp.Percentage
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

	response := CheckoutResponse{
		Products: []ProductResponse{},
	}

	for _, productRequest := range data.Products {
		productData, err := svc.Repository.GetProduct(productRequest.ID)
		if err != nil {
			log.Warnf("Failed requesting product from repository: %v", err)
			continue
		}

		productResponse := NewProductResponse(productData, productRequest)

		discountReq := &discount.GetDiscountRequest{
			ProductID: int32(productRequest.ID),
		}

		discountResp, err := svc.DiscountClient.GetDiscount(ctx, discountReq)
		if err != nil {
			log.Warnf("Failed requesting discount from external service: %v", err)
		} else {
			productResponse.Discount = CalculateDiscount(productData.Amount, discountResp.GetPercentage())
		}

		response.TotalAmount += productResponse.TotalAmount
		response.TotalDiscount += productResponse.Discount
		response.TotalAmountWithDiscount += productResponse.TotalAmount - productResponse.Discount
		response.Products = append(response.Products, *productResponse)
	}

	gift, ok := svc.BlackFridayGift()
	if ok {
		response.Products = append(response.Products, *gift)
	}

	return c.JSON(http.StatusOK, response)
}
