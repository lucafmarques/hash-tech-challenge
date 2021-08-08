package checkout

import (
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
		ID:          productData.ID,
		UnitAmount:  productData.Amount,
		TotalAmount: productData.Amount * productRequest.Quantity,
		Gift:        productData.Gift,
	}
}

func BuildIdsList(products []ProductRequest) []int {
	var ids []int

	for _, p := range products {
		ids = append(ids, p.ID)
	}

	return ids
}

func BlackFridayGift(date string, repo repository.Repository) (*ProductResponse, bool) {
	if date == time.Now().Format("01/02") {
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
		ID:          gift.ID,
		UnitAmount:  gift.Amount,
		TotalAmount: gift.Amount,
		Gift:        true,
	}, true
}
