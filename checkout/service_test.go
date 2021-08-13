package checkout

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/lucafmarques/hash-test/config"
	"gitlab.com/lucafmarques/hash-test/mocks"
)

func TestNewCheckoutService(t *testing.T) {
	config := config.ServiceConfig{
		Port:        ":TEST",
		Environment: config.DEVELOPMENT,
	}
	client := mocks.MockDiscountClient{}
	repo := mocks.MockRepository{}

	svc := NewCheckoutService(config, client, repo)

	assert.IsType(t, &Service{}, svc, "Failed asserting type of Service")
}

func TestNewCheckoutServiceDefaults(t *testing.T) {
	config := config.ServiceConfig{}
	client := mocks.MockDiscountClient{}
	repo := mocks.MockRepository{}

	svc := NewCheckoutService(config, client, repo)

	assert.IsType(t, &Service{}, svc, "Failed asserting type of Service")
}
