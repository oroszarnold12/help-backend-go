package service

import (
	"errors"
	"fmt"
	"help/dao"
	"help/dto"
	"help/errorsx"
	"help/model"
	"help/utils"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type UserService struct {
	userDao *dao.UserDao
}

func NewUserService(userDao *dao.UserDao) *UserService {
	return &UserService{userDao: userDao}
}

func (service *UserService) RegisterRoutes(authorizedRouter *mux.Router) {
	authorizedRouter.HandleFunc("/user", service.getUserByEmail).Methods("GET")
	authorizedRouter.HandleFunc("/persons", service.registerUsers).Methods("POST")
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

func (service *UserService) registerUsers(writer http.ResponseWriter, request *http.Request) {
	var userPostDtos []dto.UserPostDto
	if err := utils.ParseAndValidateSlice(request, &userPostDtos); err != nil {
		utils.WriteError(writer, errorsx.NewBadRequestError(err.Error()))
		return
	}

	existingEmails := []string{}
	for _, userPostDto := range userPostDtos {
		_, err := service.userDao.GetUserByEmail(userPostDto.Email)
		if err == nil {
			existingEmails = append(existingEmails, userPostDto.Email)
			continue
		}

		var notFoundError *errorsx.NotFoundError
		if !errors.As(err, &notFoundError) {
			utils.WriteError(writer, fmt.Errorf("Cannot check if email '%v' is used or not: %w", userPostDto.Email, err))
			return
		}
	}

	if len(existingEmails) > 0 {
		utils.WriteError(writer, errorsx.NewConflictError("User", "emails", strings.Join(existingEmails, ", ")))
		return
	}

	for _, userPostDto := range userPostDtos {
		hashedPassword, err := utils.HashPassword(userPostDto.Password)
		if err != nil {
			utils.WriteError(writer, err)
			return
		}
		userModel := model.UserModelFromPostDto(userPostDto, hashedPassword)

		if err := service.userDao.CreateUser(userModel); err != nil {
			utils.WriteError(writer, err)
			return
		}
	}

	utils.WriteJson(writer, http.StatusCreated, nil)
}
