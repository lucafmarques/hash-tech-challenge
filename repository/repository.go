package repository

import (
	_ "embed"
	"encoding/json"

	"gitlab.com/lucafmarques/hash-test/errors"
)

//go:embed products.json
var data []byte

type Repository interface {
	GetProduct(id int) (*Product, error)
	GetRandomGift() (*Product, error)
	GetAllProducts() (*[]Product, error)
}

type Embed struct {
	data []Product
}

func NewEmbedRepository() (Repository, error) {
	var products []Product
	err := json.Unmarshal(data, &products)
	if err != nil {
		return nil, errors.ErrFailedLoadingRepository(err)
	}

	return Embed{
		data: products,
	}, nil
}

func (repo Embed) GetProduct(id int) (*Product, error) {
	for i := range repo.data {
		product := repo.data[i]

		if product.Id == id {
			return &product, nil
		}
	}

	return nil, errors.ErrInvalidProductId
}

func (repo Embed) GetAllProducts() (*[]Product, error) {
	return &repo.data, nil
}

func (repo Embed) GetRandomGift() (*Product, error) {
	for i := range repo.data {
		product := repo.data[i]
		if product.Gift {
			return &product, nil
		}
	}

	return nil, errors.ErrNoGiftAvailable
}
