package checkout

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/lucafmarques/hash-test/config"
	"gitlab.com/lucafmarques/hash-test/discount"
	"gitlab.com/lucafmarques/hash-test/mocks"
	"gitlab.com/lucafmarques/hash-test/repository"
)

func TestCreateNewCore(t *testing.T) {
	config := config.CoreConfig{}
	repo := mocks.MockRepository{}
	client := mocks.MockDiscountClient{}

	core := NewCore(config, client, repo)

	assert.IsType(t, Core{}, core, "Failed asserting created core type")
}

func TestCalculateDiscountPercentage(t *testing.T) {
	ctx := context.Background()

	perc := float32(0.05)

	core := Core{
		Config: config.CoreConfig{},
		Client: mocks.MockDiscountClient{
			Resp: &discount.GetDiscountResponse{
				Percentage: perc,
			},
		},
	}

	result := core.CalculateDiscountPercentage(ctx, 0)
	assert.Equal(t, perc, result, "Failed asserting calculated discount percentage")
}

func TestCalculateDiscountPercentageNoProductWithId(t *testing.T) {
	ctx := context.Background()

	core := Core{
		Config: config.CoreConfig{},
		Client: mocks.MockDiscountClient{
			Resp: &discount.GetDiscountResponse{},
		},
	}

	result := core.CalculateDiscountPercentage(ctx, 0)
	assert.Equal(t, float32(0), result, "Failed asserting calculated discount percentage")
}

func TestCalculateCheckout(t *testing.T) {
	ctx := context.Background()

	var request []ProductRequest
	var products []*repository.Product

	perc := float32(0.05)
	ids := []int{0, 1, 2, 3, 4, 5, 6}
	sum := func(values []*repository.Product) (total int) {
		for _, v := range values {
			total += v.Amount
		}
		return
	}

	for _, id := range ids {
		request = append(request, ProductRequest{
			ID:       id,
			Quantity: 1,
		})
		products = append(products, &repository.Product{
			ID:     id,
			Amount: id * 100,
			Gift:   false,
		})
	}

	core := Core{
		Repository: mocks.MockRepository{
			Products: products,
		},
		Client: mocks.MockDiscountClient{
			Resp: &discount.GetDiscountResponse{
				Percentage: perc,
			},
		},
	}
	expectedTotalAmount := sum(products)
	expectedTotalDiscount := int(float32(sum(products)) * perc)
	resp, totalAmount, totalDiscount := core.CalculateCheckout(ctx, request)

	assert.Equal(t, len(products), len(resp), "Failed asserting that response has the same amount of products as we requested")
	assert.Equal(t, expectedTotalAmount, totalAmount, "Failed asserting that totalAmount is the same as sum of products amount")
	assert.Equal(t, expectedTotalDiscount, totalDiscount, "Failed asserting that totalDiscount is equal to totalAmount * fixed percentage")
}

func TestCalculateCheckoutNotAllRequestedProductsExist(t *testing.T) {
	ctx := context.Background()

	var request []ProductRequest
	var products []*repository.Product

	perc := float32(0.05)
	ids := []int{0, 1, 2, 3, 4, 5, 6}
	sum := func(values []*repository.Product) (total int) {
		for _, v := range values {
			total += v.Amount
		}
		return
	}

	for _, id := range ids {
		request = append(request, ProductRequest{
			ID:       id,
			Quantity: 1,
		})
		products = append(products, &repository.Product{
			ID:     id,
			Amount: id * 100,
			Gift:   false,
		})
	}

	someProducts := products[0:4]

	core := Core{
		Repository: mocks.MockRepository{
			Products: someProducts,
			Err:      errors.New("failed fetching some products"),
		},
		Client: mocks.MockDiscountClient{
			Resp: &discount.GetDiscountResponse{
				Percentage: perc,
			},
		},
	}
	expectedTotalAmount := sum(someProducts)
	expectedTotalDiscount := int(float32(sum(someProducts)) * perc)
	resp, totalAmount, totalDiscount := core.CalculateCheckout(ctx, request)

	assert.Equal(t, len(someProducts), len(resp), "Failed asserting that response has the same amount of products as we requested")
	assert.Equal(t, expectedTotalAmount, totalAmount, "Failed asserting that totalAmount is the same as sum of products amount")
	assert.Equal(t, expectedTotalDiscount, totalDiscount, "Failed asserting that totalDiscount is equal to totalAmount * fixed percentage")
}

func TestBuildIdsList(t *testing.T) {
	var products []ProductRequest

	core := Core{}
	ids := []int{0, 1, 2, 3, 4}

	for _, id := range ids {
		products = append(products, ProductRequest{
			ID: id,
		})
	}

	resp := core.BuildIdsList(products)

	assert.Equal(t, ids, resp, "Failed asserting built list of ids")
}

func TestBlackFridayGift(t *testing.T) {
	core := Core{
		Config: config.CoreConfig{
			BlackFridayDate: time.Now().Format("01/02"),
		},
		Repository: mocks.MockRepository{
			Product: &repository.Product{
				ID:   123,
				Gift: true,
			},
		},
	}

	resp, respOk := core.BlackFridayGift()

	assert.Equal(t, resp.ID, 123, "Failed asserting gift ID")
	assert.True(t, respOk, "Failed asserting gift available")
}

func TestBlackFridayGiftWrongDate(t *testing.T) {
	core := Core{
		Config: config.CoreConfig{
			BlackFridayDate: "30/02",
		},
	}

	resp, respOk := core.BlackFridayGift()

	assert.Nil(t, resp, "Failed asserting gift is nil")
	assert.False(t, respOk, "Failed asserting gift not available")
}

func TestBlackFridayGiftRepositoryFailure(t *testing.T) {
	core := Core{
		Config: config.CoreConfig{
			BlackFridayDate: time.Now().Format("01/02"),
		},
		Repository: mocks.MockRepository{
			Err: errors.New("No gift available."),
		},
	}

	resp, respOk := core.BlackFridayGift()

	assert.Nil(t, resp, "Failed asserting gift is nil")
	assert.False(t, respOk, "Failed asserting gift not available")
}
