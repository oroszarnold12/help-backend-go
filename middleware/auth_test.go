package middleware

import (
	"encoding/json"
	"fmt"
	"help/constant"
	"help/model"
	"help/test/factory"
	"help/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type userListerMock struct {
	user               *model.User
	getUserByUuidError error
}

func (userLister *userListerMock) GetUserByEmail(email string) (*model.User, error) {
	panic("not implemented")
}

func (userLister *userListerMock) GetUserByUuid(uuid uuid.UUID) (*model.User, error) {
	return userLister.user, userLister.getUserByUuidError
}

func (userLister *userListerMock) GetUsers() ([]model.User, error) {
	panic("not implemented")
}

type handlerMock struct {
	called      bool
	contextUser *model.User
}

func (handler *handlerMock) handle(writer http.ResponseWriter, request *http.Request) {
	handler.called = true
	handler.contextUser = utils.GetCurrentUser(request)
}

func TestMiddlewareFunc(t *testing.T) {
	user := factory.NewTestUser()
	token, _ := utils.CreateToken(&user)
	invalidToken := "invalid-token"
	invalidCookie := http.Cookie{
		Name:     constant.AuthCookieName,
		Value:    invalidToken,
		HttpOnly: true,
		Path:     "/",
	}
	cookie := http.Cookie{
		Name:     constant.AuthCookieName,
		Value:    token,
		HttpOnly: true,
		Path:     "/",
	}

	tests := []struct {
		name               string
		cookie             http.Cookie
		wantStatus         int
		wantError          string
		getUserByUuidError error
		allowedRoles       []model.Role
	}{
		{
			name:       "Missing cookie",
			wantStatus: http.StatusUnauthorized,
			wantError:  "Missing or invalid authorization cookie",
		},
		{
			name:       "Invalid token",
			cookie:     invalidCookie,
			wantStatus: http.StatusUnauthorized,
			wantError:  "Missing or invalid authorization cookie",
		},
		{
			name:               "Cannot get user",
			cookie:             cookie,
			wantStatus:         http.StatusUnauthorized,
			wantError:          "Missing or invalid authorization cookie",
			getUserByUuidError: fmt.Errorf("Cannot get user by UUID"),
		},
		{
			name:       "Valid cookie with valid token",
			cookie:     cookie,
			wantStatus: http.StatusOK,
		},
		{
			name:         "Has one of the allowed roles",
			cookie:       cookie,
			wantStatus:   http.StatusOK,
			allowedRoles: []model.Role{model.RoleStudent, model.RoleTeacher},
		},
		{
			name:         "Missing allowed role",
			cookie:       cookie,
			wantStatus:   http.StatusForbidden,
			wantError:    "You do not have the required role to perform this action",
			allowedRoles: []model.Role{model.RoleTeacher, model.RoleAdmin},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userListerMock := &userListerMock{user: &user, getUserByUuidError: test.getUserByUuidError}
			authMiddleware := NewAuthMiddleware(userListerMock)
			handlerMock := handlerMock{}

			router := mux.NewRouter()
			router.Use(authMiddleware.MiddlewareFunc)
			router.HandleFunc("/path", handlerMock.handle).Methods(http.MethodGet)

			if test.allowedRoles != nil {
				authMiddleware.AllowRoles("/path", http.MethodGet, test.allowedRoles)
			}

			request := httptest.NewRequest(http.MethodGet, "/path", nil)
			request.AddCookie(&test.cookie)
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)

			response := recorder.Result()
			defer response.Body.Close()

			if response.StatusCode != test.wantStatus {
				t.Errorf("Got status %d, wanted %d", response.StatusCode, test.wantStatus)
			}

			if test.wantError == "" {
				if !handlerMock.called {
					t.Errorf("Handler was not called, wanted it to be called")
				}
				if handlerMock.contextUser != &user {
					t.Errorf("Got contextUser %v, wanted %v", handlerMock.contextUser, user)
				}
			} else {
				var data map[string]string
				if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
					t.Errorf("Cannot decode body: %v", err)
				}

				error := data["error"]
				if error != test.wantError {
					t.Errorf("Got error '%v', wanted '%v'", error, test.wantError)
				}
			}
		})
	}
}
