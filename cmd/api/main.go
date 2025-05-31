package api

import (
	"database/sql"
	"fmt"
	"help/dao"
	"help/service"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Api struct {
	port int
	db   *sql.DB
}

func NewApi(port int, db *sql.DB) *Api {
	return &Api{port: port, db: db}
}

func (api *Api) Run() error {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	userDao := dao.NewUserDao(api.db)
	userService := service.NewUserService(userDao)
	userService.RegisterRoutes(subRouter)

	log.Println("Listening on port", api.port)

	return http.ListenAndServe(fmt.Sprintf(":%d", api.port), subRouter)
}
