package errorsx

import (
	"fmt"
	"net/http"
)

type ConflictError struct {
	EntityType string
	Key        string
	Value      string
}

func (error *ConflictError) Error() string {
	return fmt.Sprintf("%v with %v '%v' already exists", error.EntityType, error.Key, error.Value)
}

func (error *ConflictError) StatusCode() int {
	return http.StatusConflict
}

func NewConflictError(entityType string, key string, value string) *ConflictError {
	return &ConflictError{EntityType: entityType, Key: key, Value: value}
}
