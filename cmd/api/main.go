package api

import (
	"database/sql"
	"fmt"
	"help/dao"
	"help/middleware"
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
	publicRouter := router.PathPrefix("/api/v1").Subrouter()
	authorizedRouter := router.PathPrefix("/api/v1").Subrouter()

	userDao := dao.NewUserDao(api.db)

	loggingMiddleware := middleware.NewLoggingMiddleware()
	authMiddleware := middleware.NewAuthMiddleware(userDao)

	publicRouter.Use(loggingMiddleware.MiddlewareFunc)
	authorizedRouter.Use(loggingMiddleware.MiddlewareFunc)
	authorizedRouter.Use(authMiddleware.MiddlewareFunc)

	userService := service.NewUserService(userDao)
	userService.RegisterRoutes(authorizedRouter)

	authService := service.NewAuthService(userDao)
	authService.RegisterRoutes(publicRouter, authorizedRouter)

	log.Println("Listening on port", api.port)

	return http.ListenAndServe(fmt.Sprintf(":%d", api.port), router)
}
