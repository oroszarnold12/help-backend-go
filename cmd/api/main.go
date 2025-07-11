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
	"github.com/rs/cors"
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
	// front-end does not support api/v1 format
	// publicRouter := router.PathPrefix("/api/v1").Subrouter()
	// authorizedRouter := router.PathPrefix("/api/v1").Subrouter()

	publicRouter := router.NewRoute().Subrouter()
	authorizedRouter := router.NewRoute().Subrouter()

	userDao := dao.NewUserDao(api.db)
	courseDao := dao.NewCourseDao(api.db)
	participationDao := dao.NewPariticipationDao(api.db)
	invitationsDao := dao.NewInvitationDao(api.db)

	loggingMiddleware := middleware.NewLoggingMiddleware()
	authMiddleware := middleware.NewAuthMiddleware(userDao)

	publicRouter.Use(loggingMiddleware.MiddlewareFunc)
	authorizedRouter.Use(loggingMiddleware.MiddlewareFunc)
	authorizedRouter.Use(authMiddleware.MiddlewareFunc)

	userService := service.NewUserService(userDao)
	userService.RegisterRoutes(authMiddleware, authorizedRouter)

	authService := service.NewAuthService(userDao)
	authService.RegisterRoutes(publicRouter)

	statusService := service.NewStatusService()
	statusService.RegisterRoutes(publicRouter)

	courseService := service.NewCourseService(courseDao)
	courseService.RegisterRoutes(authMiddleware, authorizedRouter)

	participationService := service.NewParticipaionService(participationDao)
	participationService.RegisterRoutes(authorizedRouter)

	invitationsService := service.NewInvitationService(invitationsDao, participationDao)
	invitationsService.RegisterRoutes(authorizedRouter)

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8100"},
		AllowCredentials: true,
		AllowedMethods:   []string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodDelete}},
	).Handler(router)
	log.Println("Listening on port", api.port)

	return http.ListenAndServe(fmt.Sprintf(":%d", api.port), handler)
}
