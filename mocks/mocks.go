package mocks

import (
	"context"

	"gitlab.com/lucafmarques/hash-test/discount"
	"gitlab.com/lucafmarques/hash-test/repository"
	"google.golang.org/grpc"
)

type MockDiscountClient struct {
	Resp *discount.GetDiscountResponse
	Err  error
}

func (d MockDiscountClient) GetDiscount(ctx context.Context, in *discount.GetDiscountRequest, opts ...grpc.CallOption) (*discount.GetDiscountResponse, error) {
	return d.Resp, d.Err
}

type MockRepository struct {
	Product  *repository.Product
	Products []*repository.Product
	Err      error
}

func (r MockRepository) GetProduct(id int) (*repository.Product, error) {
	return r.Product, r.Err
}

func (r MockRepository) GetRandomGift() (*repository.Product, error) {
	return r.Product, r.Err
}

func (r MockRepository) GetProductsByIds(ids []int) (map[int]*repository.Product, error) {
	m := make(map[int]*repository.Product)

	for _, product := range r.Products {
		m[product.ID] = product
	}

	return m, r.Err
}
