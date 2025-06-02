package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func parseJSON(request *http.Request, target any) error {
	if request.Body == nil {
		return fmt.Errorf("Missing request body")
	}

	if err := json.NewDecoder(request.Body).Decode(target); err != nil {
		return fmt.Errorf("Cannot decode body as JSON")
	}

	return nil
}

func GetCookieFromRequest(request *http.Request, cookieName string) (*http.Cookie, error) {
	cookie, err := request.Cookie(cookieName)
	if err != nil {
		return nil, fmt.Errorf("Cannot get cookie '%v' from request: %w", cookieName, err)
	}

	return cookie, nil
}

func ParseAndValidateSlice(request *http.Request, target any) error {
	if err := parseJSON(request, &target); err != nil {
		return err
	}

	if err := Validator.Var(target, "dive"); err != nil {
		return err
	}

	return nil
}

func ParseAndValidateStruct(request *http.Request, target any) error {
	if err := parseJSON(request, &target); err != nil {
		return err
	}

	if err := Validator.Struct(target); err != nil {
		return err
	}

	return nil
}
