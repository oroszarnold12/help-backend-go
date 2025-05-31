package errorsx

import (
	"fmt"
	"net/http"
)

type BadRequestError struct {
	Reason string
}

func (error *BadRequestError) Error() string {
	return fmt.Sprintf("Bad request: %s", error.Reason)
}

func (error *BadRequestError) StatusCode() int {
	return http.StatusBadRequest
}

func NewBadRequestError(reason string) *BadRequestError {
	return &BadRequestError{Reason: reason}
}
