package errorsx

import (
	"net/http"
)

type ForbiddenError struct{}

func (error *ForbiddenError) Error() string {
	return "You do not have the required role to perform this action"
}

func (error *ForbiddenError) StatusCode() int {
	return http.StatusForbidden
}

func NewForbiddenError() *ForbiddenError {
	return &ForbiddenError{}
}
