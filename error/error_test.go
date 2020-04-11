package error

import (
	"errors"
	"testing"
)

func TestMyError_Error(t *testing.T) {
	err := errors.New("test error")

	myErr := WrapError(err, "test error is wrapped")
	HandleError(1, myErr)
}
