package repository

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/lucafmarques/hash-test/config"
)

func TestNewEmbedRepository(t *testing.T) {
	host := "TEST"
	config := config.RepositoryConfig{Host: host}

	repo, err := NewEmbedRepository(config)

	assert.Equal(t, host, repo.(Embed).Config.Host, "Failed asserting same host in repo config")
	assert.Nil(t, err, "Failed asserting nil err")
}
func TestNewEmbedRepositoryError(t *testing.T) {
	host := "TEST"
	data = []byte("invalid_data")
	config := config.RepositoryConfig{Host: host}

	repo, err := NewEmbedRepository(config)

	assert.Nil(t, repo, "Failed asserting nil repo")
	assert.Error(t, err, "Failed asserting err")
}

func TestGetProductsByIds(t *testing.T) {
	data := map[int]Product{}
	ids := []int{0, 1, 2, 3, 4}

	for _, id := range ids {
		data[id] = Product{
			ID: id,
		}
	}

	repo := Embed{
		data:   data,
		Config: config.RepositoryConfig{},
	}

	products, err := repo.GetProductsByIds(ids)

	assert.IsType(t, &Product{}, products[0], "Failed asserting Product type")
	assert.Nil(t, err, "Failed asserting nil err")
}

func TestGetProductsByIdsNotAllProductsExist(t *testing.T) {
	data := map[int]Product{}
	ids := []int{0, 1, 2, 3, 4}

	for _, id := range ids[:3] {
		data[id] = Product{
			ID: id,
		}
	}

	repo := Embed{
		data:   data,
		Config: config.RepositoryConfig{},
	}

	products, err := repo.GetProductsByIds(ids)

	assert.IsType(t, &Product{}, products[0], "Failed asserting Product type")
	assert.Error(t, err, "Failed asserting err")
	assert.Contains(t, err.Error(), "ID=4", "Failed asserting err with no missing product string")

	_, ok := products[4]
	assert.False(t, ok, "Failed asserting false ok accessing product by ID")
}

func TestGetProduct(t *testing.T) {
	id := 0
	data := map[int]Product{
		id: {
			ID: id,
		},
	}

	repo := Embed{
		data:   data,
		Config: config.RepositoryConfig{},
	}

	product, err := repo.GetProduct(id)

	assert.IsType(t, &Product{}, product, "Failed asserting Product type")
	assert.Nil(t, err, "Failed asserting nil err")
}

func TestGetProductInvalidProductId(t *testing.T) {
	id := 0
	data := map[int]Product{}

	repo := Embed{
		data:   data,
		Config: config.RepositoryConfig{},
	}

	product, err := repo.GetProduct(id)

	assert.Nil(t, product, "Failed asserting nil")
	assert.Error(t, err, "Failed asserting err")
	assert.Contains(t, err.Error(), fmt.Sprintf("ID=%v", id), "Failed asserting err with no missing product string")
}

func TestGetRandomGift(t *testing.T) {
	id := 0
	data := map[int]Product{
		id: {
			ID:   id,
			Gift: true,
		},
	}

	repo := Embed{
		data:   data,
		Config: config.RepositoryConfig{},
	}

	product, err := repo.GetRandomGift()

	assert.IsType(t, &Product{}, product, "Failed asserting Product type")
	assert.Nil(t, err, "Failed asserting nil err")
}

func TestGetRandomGiftNoGiftAvailable(t *testing.T) {
	id := 0
	data := map[int]Product{
		id: {
			ID:   id,
			Gift: false,
		},
	}

	repo := Embed{
		data:   data,
		Config: config.RepositoryConfig{},
	}

	product, err := repo.GetRandomGift()

	assert.Nil(t, product, "Failed asserting nil")
	assert.Error(t, err, "Failed asserting err")
}
