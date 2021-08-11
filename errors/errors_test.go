package errors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrInvalidProductId(t *testing.T) {
	id := 20

	err := ErrInvalidProductId(id)
	assert.Error(t, err, "Failed asserting return was an error")
	assert.Contains(t, err.Error(), fmt.Sprint(id), "Failed asserting error contains correct ID")
}

func TestFailedLoadingRepository(t *testing.T) {
	argErr := errors.New("Test error")

	err := ErrFailedLoadingRepository(argErr)
	assert.Error(t, err, "Failed asserting return was an error")
	assert.ErrorIs(t, err, argErr, "Failed asserting error has argError in chain")
}
