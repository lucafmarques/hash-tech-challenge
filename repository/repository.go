package repository

import (
	_ "embed"
	"encoding/json"

	"github.com/labstack/gommon/log"
	"gitlab.com/lucafmarques/hash-test/config"
	"gitlab.com/lucafmarques/hash-test/errors"
)

//go:embed products.json
var data []byte

type Repository interface {
	GetProduct(id int) (*Product, error)
	GetRandomGift() (*Product, error)
	GetAllProducts() (*[]Product, error)
	GetProductsById(ids []int) map[int]*Product
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

func (repo Embed) GetProductsById(ids []int) map[int]*Product {
	products := map[int]*Product{}

	for _, id := range ids {
		product, ok := repo.data[id]
		if !ok {
			log.Infof("failed fetching data for product with ID=%v", id)
			continue
		}

		products[id] = &product
	}

	return products
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
	for _, p := range repo.data {
		if p.Gift {
			return &p, nil
		}
	}

	return nil, errors.ErrNoGiftAvailable
}
