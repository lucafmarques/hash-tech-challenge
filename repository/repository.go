package repository

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"gitlab.com/lucafmarques/hash-test/config"
	"gitlab.com/lucafmarques/hash-test/errors"
)

//go:embed products.json
var data []byte

type Repository interface {
	GetProduct(id int) (*Product, error)
	GetRandomGift() (*Product, error)
	GetAllProducts() (*[]Product, error)
	GetProductsByIds(ids []int) (map[int]*Product, error)
}

type Embed struct {
	data   map[int]Product
	Config config.RepositoryConfig
}

func NewEmbedRepository(config config.RepositoryConfig) (Repository, error) {
	var products []Product
	err := json.Unmarshal(data, &products)
	if err != nil {
		return nil, errors.ErrFailedLoadingRepository(err)
	}

	productMap := map[int]Product{}

	for _, p := range products {
		productMap[p.ID] = p
	}

	return Embed{
		data:   productMap,
		Config: config,
	}, nil
}

func (repo Embed) GetProductsByIds(ids []int) (map[int]*Product, error) {
	var errors []string
	products := map[int]*Product{}

	for _, id := range ids {
		product, err := repo.GetProduct(id)
		if err != nil {
			errors = append(errors, err.Error())
			continue
		}

		products[id] = product
	}

	if len(errors) > 0 {
		return products, fmt.Errorf("failed getting some products: %v", strings.Join(errors, ", "))
	}

	return products, nil
}

func (repo Embed) GetProduct(id int) (*Product, error) {
	product, ok := repo.data[id]
	if !ok {
		return nil, errors.ErrInvalidProductId(id)
	}

	return &product, nil
}

func (repo Embed) GetAllProducts() (*[]Product, error) {
	var products []Product

	for _, p := range repo.data {
		products = append(products, p)
	}

	return &products, nil
}

func (repo Embed) GetRandomGift() (*Product, error) {
	// This implementation is semi-random, because for Embed's Repository
	// implementation, data is stored in a map which has semi-random
	// access, not preserving order when being iterated over.
	for _, p := range repo.data {
		if p.Gift {
			return &p, nil
		}
	}

	return nil, errors.ErrNoGiftAvailable
}
