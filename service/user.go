package service

import (
	"help/dao"
	"help/errorsx"
	"help/utils"
	"net/http"

	"github.com/gorilla/mux"
)

type UserService struct {
	userDao *dao.UserDao
}

func NewUserService(userDao *dao.UserDao) *UserService {
	return &UserService{userDao: userDao}
}

func (service *UserService) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/user", service.getUserByEmail).Methods("GET")
}

func (service *UserService) getUserByEmail(writer http.ResponseWriter, request *http.Request) {
	if !request.URL.Query().Has("email") {
		utils.WriteError(writer, errorsx.NewBadRequestError("Missing email query parameter"))
		return
	}
	email := request.URL.Query().Get("email")

	user, err := service.userDao.GetUserByEmail(email)
	if err != nil {
		utils.WriteError(writer, err)
		return
	}

	utils.WriteJson(writer, http.StatusOK, user.ToDto())
}
