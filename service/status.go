package service

import (
	"help/utils"
	"net/http"

	"github.com/gorilla/mux"
)

type StatusService struct{}

func NewStatusService() *StatusService {
	return &StatusService{}
}

func (service *StatusService) RegisterRoutes(publicRotuer *mux.Router) {
	publicRotuer.HandleFunc("/server-status/ping", service.ping).Methods(http.MethodGet)
}

func (service *StatusService) ping(writer http.ResponseWriter, request *http.Request) {
	utils.WriteJson(writer, http.StatusOK, nil)
}
