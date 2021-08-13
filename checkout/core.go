package checkout

import (
	"context"
	"sync"
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

func NewCore(config config.CoreConfig, client discount.DiscountClient, repo repository.Repository) Core {
	return Core{
		Config:     config,
		Repository: repo,
		Client:     client,
	}
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

func (c Core) CalculateCheckout(ctx context.Context, requestedProducts []ProductRequest) ([]ProductResponse, int, int, bool) {
	var (
		totalAmount   int
		totalDiscount int
		hasGift       bool
		wg            sync.WaitGroup
	)

	response := []ProductResponse{}
	ch := make(chan ProductResponse, len(requestedProducts))

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
		if product.Gift && hasGift {
			log.Info("Checkout car already has a gift, skipping product")
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

		if product.Gift {
			hasGift = true
		}

		wg.Add(1)
		go func(id int, p ProductResponse, ctx context.Context) {
			p.Discount = int(float32(p.UnitAmount) * c.CalculateDiscountPercentage(ctx, id))
			ch <- p
			wg.Done()
		}(r.ID, *productResponse, ctx)
	}
	wg.Wait()
	close(ch)

	for p := range ch {
		totalAmount += p.TotalAmount
		totalDiscount += p.Discount
		response = append(response, p)
	}

	return response, totalAmount, totalDiscount, hasGift
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
