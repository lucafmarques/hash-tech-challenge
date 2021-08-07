package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gitlab.com/lucafmarques/hash-test/discount"
)

var Data []Product

func init() {
	json.Unmarshal(Products, &Data)
}

type CheckoutResponse struct {
	TotalAmount             int               `json:"total_amount"`
	TotalAmountWithDiscount int               `json:"total_amount_with_discount"`
	TotalDiscount           int               `json:"total_discount"`
	Products                []ProductResponse `json:"products"`
}

type ProductResponse struct {
	Id          int  `json:"id"`
	Quantity    int  `json:"quantity"`
	UnitAmount  int  `json:"unit_amount"`
	TotalAmount int  `json:"total_amount"`
	Discount    int  `json:"discount"`
	Gift        bool `json:"is_gift"`
}

type DiscountResponse struct {
	Percentage float32 `json:"percentage"`
}

type Product struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Amount      int    `json:"amount"`
	Gift        bool   `json:"is_gift"`
}

func (svc *Service) HandleHello(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "Hello World!"})
}

func (svc *Service) GetAllProducts(c echo.Context) error {
	return c.JSON(http.StatusOK, Data)
}

func (svc *Service) GetProductDiscount(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.Atoi((c.Param("id")))
	if err != nil {
		svc.Server.Logger.Warnf("Failed conversion of Product id to int: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Product ID should be a number.")
	}

	resp := DiscountResponse{
		Percentage: 0,
	}

	req := &discount.GetDiscountRequest{
		ProductID: int32(id),
	}

	discountResp, err := svc.DiscountClient.GetDiscount(ctx, req)
	if err != nil {
		svc.Server.Logger.Warnf("Failed requesting discount from external service: %v", err)
	} else {
		resp.Percentage = discountResp.Percentage
	}

	return c.JSON(http.StatusOK, resp)
}
