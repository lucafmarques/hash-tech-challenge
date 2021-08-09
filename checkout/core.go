package checkout

import (
	"context"
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
	"gitlab.com/lucafmarques/hash-test/config"
	"gitlab.com/lucafmarques/hash-test/discount"
	"gitlab.com/lucafmarques/hash-test/repository"
)

type Core struct {
	Config     config.CoreConfig
	Repository repository.Repository
	Client     discount.DiscountClient
}

func (c Core) AllProducts() (*[]repository.Product, error) {
	products, err := c.Repository.GetAllProducts()
	if err != nil {
		return nil, fmt.Errorf("failed fetching products: %v", err.Error())
	}

	return products, nil
}

func (c Core) CalculateDiscountPercentage(ctx context.Context, id int) float32 {
	var percentage float32

	req := &discount.GetDiscountRequest{
		ProductID: int32(id),
	}

	discountResp, err := c.Client.GetDiscount(ctx, req)
	if err != nil {
		log.Warnf("Failed requesting discount from external service: %v", err.Error())
	} else {
		percentage = discountResp.GetPercentage()
	}

	return percentage
}

func (c Core) CalculateCheckout(ctx context.Context, requestedProducts []ProductRequest) ([]ProductResponse, int, int) {
	var (
		response      []ProductResponse
		totalAmount   int
		totalDiscount int
	)

	products, err := c.Repository.GetProductsByIds(c.BuildIdsList(requestedProducts))
	if err != nil {
		// Not exiting here with an error allows a bad agent to probe
		// the database for inexistent ID's.
		// This could be changed to account for that.
		log.Warnf("Failed getting all requested products: %v", err)
	}

	for _, r := range requestedProducts {
		product, ok := products[r.ID]
		if !ok {
			log.Warnf("No product data exists for ID=%v", r.ID)
			continue
		}

		productResponse := &ProductResponse{
			Discount:    0,
			Quantity:    r.Quantity,
			ID:          product.ID,
			UnitAmount:  product.Amount,
			TotalAmount: product.Amount * r.Quantity,
			Gift:        product.Gift,
		}

		percentage := c.CalculateDiscountPercentage(ctx, r.ID)
		productResponse.Discount = int(float32(product.Amount) * percentage)

		totalAmount += productResponse.TotalAmount
		totalDiscount += productResponse.Discount
		response = append(response, *productResponse)
	}

	return response, totalAmount, totalDiscount
}

func (c Core) BuildIdsList(products []ProductRequest) []int {
	var ids []int

	for _, p := range products {
		ids = append(ids, p.ID)
	}

	return ids
}

func (c Core) BlackFridayGift() (*ProductResponse, bool) {
	if time.Now().Format("01/02") != c.Config.BlackFridayDate {
		return nil, false
	}

	gift, err := c.Repository.GetRandomGift()
	if err != nil {
		log.Warnf("Failed requesting gift from repository: %v", err.Error())
		return nil, false
	}

	return &ProductResponse{
		Discount:    0,
		Quantity:    1,
		ID:          gift.ID,
		UnitAmount:  0,
		TotalAmount: 0,
		Gift:        gift.Gift,
	}, true
}
