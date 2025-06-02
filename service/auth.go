package service

import (
	"errors"
	"fmt"
	"help/constant"
	"help/dao"
	"help/dto"
	"help/errorsx"
	"help/utils"
	"net/http"

	"github.com/gorilla/mux"
)

type AuthService struct {
	userDao *dao.UserDao
}

func NewAuthService(userDao *dao.UserDao) *AuthService {
	return &AuthService{userDao: userDao}
}

func (service *AuthService) RegisterRoutes(publicRouter *mux.Router, authorizedRouter *mux.Router) {
	publicRouter.HandleFunc("/auth/login", service.login).Methods("POST")
	authorizedRouter.HandleFunc("/auth/logout", service.logout).Methods("GET")
}

func (service *AuthService) login(writer http.ResponseWriter, request *http.Request) {
	var loginDto dto.LoginDto
	if err := utils.ParseAndValidateStruct(request, &loginDto); err != nil {
		utils.WriteError(writer, errorsx.NewBadRequestError(err.Error()))
		return
	}

	user, err := service.userDao.GetUserByEmail(loginDto.Username)
	if err != nil {
		var notFoundError *errorsx.NotFoundError
		if errors.As(err, &notFoundError) {
			utils.WriteError(writer, errorsx.NewBadRequestError("Incorrect username or password"))
		} else {
			utils.WriteError(writer, fmt.Errorf("Cannot get user by email '%v'", loginDto.Username))
		}

		return
	}

	if !utils.ComparePasswords(user.Password, loginDto.Password) {
		utils.WriteError(writer, errorsx.NewBadRequestError("Incorrect username or password"))
		return
	}

	token, err := utils.CreateToken(user)
	if err != nil {
		utils.WriteError(writer, err)
	}

	cookie := http.Cookie{
		Name:     constant.AuthCookieName,
		Value:    token,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(writer, &cookie)
	utils.WriteJson(writer, http.StatusOK, user.ToDto())
}

func (service *AuthService) logout(writer http.ResponseWriter, request *http.Request) {
	cookie, err := utils.GetCookieFromRequest(request, constant.AuthCookieName)
	if err != nil {
		utils.WriteError(writer, fmt.Errorf("%w: %v", errorsx.NewUnauthorizedError(), err))
		return
	}
	cookie.Path = "/"
	cookie.MaxAge = -1

	http.SetCookie(writer, cookie)
}
