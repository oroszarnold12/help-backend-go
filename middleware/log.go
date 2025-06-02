package middleware

import (
	"log"
	"net/http"
)

type LoggingMiddleware struct{}

func NewLoggingMiddleware() *LoggingMiddleware {
	return &LoggingMiddleware{}
}

func (middleware *LoggingMiddleware) MiddlewareFunc(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("%s %s", request.Method, request.URL)
		handler.ServeHTTP(writer, request)
	})
}
