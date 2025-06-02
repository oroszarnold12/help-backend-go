package middleware

import (
	"context"
	"fmt"
	"help/constant"
	"help/dao"
	"help/errorsx"
	"help/utils"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	userDao *dao.UserDao
}

func NewAuthMiddleware(userDao *dao.UserDao) *AuthMiddleware {
	return &AuthMiddleware{userDao: userDao}
}

func (middleware *AuthMiddleware) MiddlewareFunc(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		cookie, err := utils.GetCookieFromRequest(request, constant.AuthCookieName)
		if err != nil {
			utils.WriteError(writer, fmt.Errorf("%w: %v", errorsx.NewUnauthorizedError(), err))
			return
		}

		token, err := utils.ValidateToken(cookie.Value)
		if err != nil {
			utils.WriteError(writer, fmt.Errorf("%w: %v", errorsx.NewUnauthorizedError(), err))
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userUuid := claims[constant.UserUuidClaimKey].(string)

		user, err := middleware.userDao.GetUserByUuid(userUuid)
		if err != nil {
			utils.WriteError(writer, fmt.Errorf("%w: %v", errorsx.NewUnauthorizedError(), err))
			return
		}

		ctx := request.Context()
		ctx = context.WithValue(ctx, constant.UserContextKey, user)
		request = request.WithContext(ctx)

		handler.ServeHTTP(writer, request)
	})
}
