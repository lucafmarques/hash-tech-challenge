package checkout

import (
	"os"
	"time"

	"github.com/labstack/gommon/log"
	"gitlab.com/lucafmarques/hash-test/repository"
)

func CalculateDiscount(productValue int, discountPerc float32) int {
	return int(float32(productValue) * discountPerc)
}

func NewProductResponse(productData *repository.Product, productRequest ProductRequest) *ProductResponse {
	return &ProductResponse{
		Discount:    0,
		Quantity:    productRequest.Quantity,
		ID:          productData.Id,
		UnitAmount:  productData.Amount,
		TotalAmount: productData.Amount * productRequest.Quantity,
		Gift:        productData.Gift,
	}
}

func BlackFridayGift(repo repository.Repository) (*ProductResponse, bool) {
	date := os.Getenv("BLACK_FRIDAY_DATE")

	ok := date == time.Now().Format("2006-01-02")
	if !ok {
		return nil, false
	}

	gift, err := repo.GetRandomGift()
	if err != nil {
		log.Warnf("Failed requesting gift from repository: %v", err)
		return nil, false
	}

	return &ProductResponse{
		Discount:    gift.Amount,
		Quantity:    1,
		ID:          gift.Id,
		UnitAmount:  gift.Amount,
		TotalAmount: gift.Amount,
		Gift:        true,
	}, true
}
