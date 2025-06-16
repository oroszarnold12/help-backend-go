package service

import (
	"errors"
	"fmt"
	"help/dao"
	"help/dto"
	"help/errorsx"
	"help/middleware"
	"help/model"
	"help/utils"
	"net/http"
	"slices"
	"strings"

	"github.com/gorilla/mux"
)

type UserService struct {
	userDao *dao.UserDao
}

func NewUserService(userDao *dao.UserDao) *UserService {
	return &UserService{userDao: userDao}
}

func (service *UserService) RegisterRoutes(authMiddleware *middleware.AuthMiddleware, authorizedRouter *mux.Router) {
	authorizedRouter.HandleFunc("/user", service.getCurrentUser).Methods(http.MethodGet)

	authorizedRouter.HandleFunc("/persons", service.registerUsers).Methods(http.MethodPost)
	authMiddleware.AllowRoles("/persons", http.MethodPost, []model.Role{model.RoleAdmin})

	authorizedRouter.HandleFunc("/persons", service.getUsers).Methods(http.MethodGet)
}

func (service *UserService) getCurrentUser(writer http.ResponseWriter, request *http.Request) {
	user := utils.GetCurrentUser(request)

	utils.WriteJson(writer, http.StatusOK, user.ToDto())
}

func (service *UserService) getUsers(writer http.ResponseWriter, request *http.Request) {
	currentUser := utils.GetCurrentUser(request)
	isAdmin := currentUser.Role == model.RoleAdmin

	users, err := service.userDao.GetUsers()
	if err != nil {
		utils.WriteError(writer, err)
	}

	if !isAdmin {
		users = slices.DeleteFunc(users, func(user model.User) bool {
			return user.Role == model.RoleAdmin
		})
	}

	utils.WriteJson(writer, http.StatusOK, map[string][]dto.UserGetDto{"persons": dto.ModelsToDtos(users)})
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
