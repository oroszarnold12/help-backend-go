package utils

import (
	"encoding/json"
	"errors"
	"help/errorsx"
	"log"
	"net/http"
)

func WriteJson(writer http.ResponseWriter, status int, data any) {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(status)

	if data != nil {
		err := json.NewEncoder(writer).Encode(data)

		if err != nil {
			log.Printf("Failed to encode %v to JSON, %v", data, err)
			writer.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func WriteError(writer http.ResponseWriter, error error) {
	log.Printf("%T: %[1]v", error)

	var statusCoder errorsx.StatusCoder
	if errors.As(error, &statusCoder) {
		WriteJson(writer, statusCoder.StatusCode(), map[string]string{"error": error.Error()})
	} else {
		WriteJson(writer, http.StatusInternalServerError, map[string]string{"error": "Internal server error, please try again later"})
	}
}
