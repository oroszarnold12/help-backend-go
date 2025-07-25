package middleware

import (
	"context"
	"fmt"
	"help/constant"
	"help/dao"
	"help/errorsx"
	"help/model"
	"help/utils"
	"net/http"
	"slices"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthMiddleware struct {
	userDao      *dao.UserDao
	allowedRoles map[string][]model.Role
}

func NewAuthMiddleware(userDao *dao.UserDao) *AuthMiddleware {
	return &AuthMiddleware{userDao: userDao, allowedRoles: make(map[string][]model.Role)}
}

func (middleware *AuthMiddleware) AllowRoles(path string, method string, roles []model.Role) {
	key := allowedRolesKey(path, method)
	if _, ok := middleware.allowedRoles[key]; ok {
		panic(fmt.Sprintf("Cannot set allowed roles to the same path and method '%s' multiple times", key))
	}

	middleware.allowedRoles[key] = roles
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

		user, err := middleware.userDao.GetUserByUuid(uuid.MustParse(userUuid))
		if err != nil {
			utils.WriteError(writer, fmt.Errorf("%w: %v", errorsx.NewUnauthorizedError(), err))
			return
		}

		if roles, ok := middleware.allowedRoles[allowedRolesKey(request.URL.Path, request.Method)]; ok {
			if !slices.Contains(roles, user.Role) {
				utils.WriteError(writer, fmt.Errorf("%w: User '%v' does not have any of the allowed roles: %v", errorsx.NewForbiddenError(), user.Email, roles))
				return
			}
		}

		ctx := request.Context()
		ctx = context.WithValue(ctx, constant.UserContextKey, user)
		request = request.WithContext(ctx)

		handler.ServeHTTP(writer, request)
	})
}

func allowedRolesKey(path string, method string) string {
	return fmt.Sprintf("%s: %s", method, path)
}
