package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ParseJSON(request *http.Request, target any) error {
	if request.Body == nil {
		return fmt.Errorf("Missing request body")
	}

	if err := json.NewDecoder(request.Body).Decode(target); err != nil {
		return fmt.Errorf("Cannot decode body as JSON")
	}

	return nil
}
