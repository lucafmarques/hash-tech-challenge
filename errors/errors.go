package errors

import (
	"errors"
	"fmt"
)

var (
	ErrNoGiftAvailable = errors.New("no random gift available")
)

func ErrInvalidProductId(id int) error {
	return fmt.Errorf("no product with ID=%v exists", id)
}

func ErrFailedLoadingRepository(err error) error {
	return fmt.Errorf("failed loading data into repository: %w", err)
}
