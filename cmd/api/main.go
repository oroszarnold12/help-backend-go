package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Api struct {
	port int
}

func NewApi(port int) *Api {
	return &Api{port: port}
}

func (api *Api) Run() error {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	log.Println("Listening on port", api.port)

	return http.ListenAndServe(fmt.Sprintf(":%d", api.port), subRouter)
}
