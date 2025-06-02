package errorsx

import (
	"net/http"
)

type UnauthorizedError struct{}

func (error *UnauthorizedError) Error() string {
	return "Missing or invalid authorization cookie"
}

func (error *UnauthorizedError) StatusCode() int {
	return http.StatusUnauthorized
}

func NewUnauthorizedError() *UnauthorizedError {
	return &UnauthorizedError{}
}
