package service

import (
	"errors"
	"fmt"
	"help/constant"
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
	authorizedRouter.HandleFunc("/user", service.getCurrentUser).Methods("GET")
	authorizedRouter.HandleFunc("/persons", service.registerUsers).Methods("POST")
	authorizedRouter.HandleFunc("/persons", service.getUsers).Methods("GET")
}

func (service *UserService) getCurrentUser(writer http.ResponseWriter, request *http.Request) {
	user := request.Context().Value(constant.UserContextKey).(*model.User)

	utils.WriteJson(writer, http.StatusOK, user.ToDto())
}

func (service *UserService) getUsers(writer http.ResponseWriter, request *http.Request) {
	currentUser := request.Context().Value(constant.UserContextKey).(*model.User)
	isAdmin := currentUser.Role == model.RoleAdmin

	users, err := service.userDao.GetUsers()
	if err != nil {
		utils.WriteError(writer, err)
	}

	if isAdmin {
		userDtos := make([]dto.UserGetDto, len(users))
		for index, user := range users {
			userDtos[index] = user.ToDto()
		}

		utils.WriteJson(writer, http.StatusOK, map[string][]dto.UserGetDto{"persons": userDtos})
	} else {
		var nonAdminUserDtos []dto.UserGetDto

		for _, user := range users {
			if user.Role != model.RoleAdmin {
				nonAdminUserDtos = append(nonAdminUserDtos, user.ToDto())
			}
		}

		utils.WriteJson(writer, http.StatusOK, map[string][]dto.UserGetDto{"persons": nonAdminUserDtos})
	}
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
		user := model.UserFromPostDto(userPostDto, hashedPassword)

		if err := service.userDao.CreateUser(user); err != nil {
			utils.WriteError(writer, err)
			return
		}
	}

	utils.WriteJson(writer, http.StatusCreated, nil)
}
