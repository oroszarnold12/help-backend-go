package errorsx

import (
	"fmt"
	"net/http"
)

type NotFoundError struct {
	EntityType string
	Key        string
}

func (error *NotFoundError) Error() string {
	return fmt.Sprintf("%v not found by key %v", error.EntityType, error.Key)
}

func (error *NotFoundError) StatusCode() int {
	return http.StatusNotFound
}

func NewNotFoundError(entityType string, key string) *NotFoundError {
	return &NotFoundError{EntityType: entityType, Key: key}
}
