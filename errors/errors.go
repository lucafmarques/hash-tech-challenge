package errors

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidProductId = errors.New("no product with that ID exists")
	ErrNoGiftAvailable  = errors.New("no random gift available")
)

func ErrFailedLoadingRepository(err error) error {
	return fmt.Errorf("failed loading data into repository: %w", err)
}
