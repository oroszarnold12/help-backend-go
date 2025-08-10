package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"help/constant"
	"help/dao"
	"help/dto"
	"help/errorsx"
	"help/middleware"
	"help/model"
	"help/test/factory"
	"help/utils"
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type userListSaverMock struct {
	users               []model.User
	userByEmailCalls    []string
	userByUuidCalls     []uuid.UUID
	createUserCalls     []model.User
	getUsersError       error
	getUserByEmailError error
	getUserByUuidError  error
	createUserError     error
}

func (listSaver *userListSaverMock) CreateUser(user model.User) error {
	listSaver.createUserCalls = append(listSaver.createUserCalls, user)
	if listSaver.createUserError != nil {
		return listSaver.createUserError
	}
	listSaver.users = append(listSaver.users, user)
	return nil
}

func (listSaver *userListSaverMock) GetUserByEmail(email string) (*model.User, error) {
	listSaver.userByEmailCalls = append(listSaver.userByEmailCalls, email)
	if listSaver.getUserByEmailError != nil {
		return nil, listSaver.getUserByEmailError
	}
	index := slices.IndexFunc(listSaver.users, func(user model.User) bool { return user.Email == email })
	if index == -1 {
		return nil, errorsx.NewNotFoundError("User", email)
	}

	return &listSaver.users[index], nil
}

func (listSaver *userListSaverMock) GetUserByUuid(uuid uuid.UUID) (*model.User, error) {
	listSaver.userByUuidCalls = append(listSaver.userByUuidCalls, uuid)
	if listSaver.getUserByUuidError != nil {
		return nil, listSaver.getUserByUuidError
	}
	index := slices.IndexFunc(listSaver.users, func(user model.User) bool { return user.Uuid == uuid })
	if index == -1 {
		return nil, errorsx.NewNotFoundError("User", uuid.String())
	}

	return &listSaver.users[index], nil
}

func (listSaver *userListSaverMock) GetUsers() ([]model.User, error) {
	return listSaver.users, listSaver.getUsersError
}

func newTestRouter(t *testing.T, userListSaver dao.UserListSaver, preAuthorizer middleware.RoutePreAuthorizer) *mux.Router {
	t.Helper()

	userService := NewUserService(userListSaver)
	router := mux.NewRouter()
	userService.RegisterRoutes(preAuthorizer, router)

	return router
}

type allowRolesCall struct {
	Path   string
	Method string
	Roles  []model.Role
}

type preAuthorizerMock struct {
	allowRolesCalls []allowRolesCall
}

func (preAuthorizer *preAuthorizerMock) AllowRoles(path, method string, roles []model.Role) {
	preAuthorizer.allowRolesCalls = append(preAuthorizer.allowRolesCalls, allowRolesCall{Path: path, Method: method, Roles: roles})
}

func TestRegisterRoutes(t *testing.T) {
	preAuthorizerMock := &preAuthorizerMock{}
	userListSaverMock := &userListSaverMock{}
	newTestRouter(t, userListSaverMock, preAuthorizerMock)

	wantAllowRolesCalls := []allowRolesCall{{Path: "/persons", Method: http.MethodPost, Roles: []model.Role{model.RoleAdmin}}}

	if !cmp.Equal(wantAllowRolesCalls, preAuthorizerMock.allowRolesCalls) {
		t.Errorf("got requireRolesCalls %v, wanted %v", preAuthorizerMock.allowRolesCalls, wantAllowRolesCalls)
	}
}

func TestGetUsers(t *testing.T) {
	currentUser := factory.NewTestUser()
	currentAdminUser := factory.NewTestUser(func(u *model.User) { u.Role = model.RoleAdmin })
	user1 := factory.NewTestUser()
	user2 := factory.NewTestUser()
	user3 := factory.NewTestUser()
	adminUser := factory.NewTestUser(func(u *model.User) { u.Role = model.RoleAdmin })

	tests := []struct {
		name        string
		users       []model.User
		error       error
		currentUser model.User
		wantStatus  int
		wantError   string
		wantUsers   []dto.UserGetDto
	}{
		{
			name:        "No users",
			users:       []model.User{},
			currentUser: currentUser,
			wantStatus:  http.StatusOK,
			wantUsers:   []dto.UserGetDto{},
		},
		{
			name:        "Returns all users",
			users:       []model.User{user1, user2, user3},
			currentUser: currentUser,
			wantStatus:  http.StatusOK,
			wantUsers:   dto.ModelsToDtos([]model.User{user1, user2, user3}),
		},
		{
			name:        "Filters out admin user for non-admin user",
			users:       []model.User{user1, user2, user3, adminUser},
			currentUser: currentUser,
			wantStatus:  http.StatusOK,
			wantUsers:   dto.ModelsToDtos([]model.User{user1, user2, user3}),
		},
		{
			name:        "Returns admin user for admin user",
			users:       []model.User{user1, user2, user3, adminUser},
			currentUser: currentAdminUser,
			wantStatus:  http.StatusOK,
			wantUsers:   dto.ModelsToDtos([]model.User{user1, user2, user3, adminUser}),
		},
		{
			name:        "Error getting users",
			error:       fmt.Errorf("Cannot get users"),
			currentUser: currentUser,
			wantStatus:  http.StatusInternalServerError,
			wantError:   "Internal server error, please try again later",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userListSaverMock := &userListSaverMock{users: test.users, getUsersError: test.error}
			preAuthorizerMock := &preAuthorizerMock{}
			router := newTestRouter(t, userListSaverMock, preAuthorizerMock)

			request := httptest.NewRequest(http.MethodGet, "/persons", nil)
			request = request.WithContext(context.WithValue(request.Context(), constant.UserContextKey, &test.currentUser))
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)

			response := recorder.Result()
			defer response.Body.Close()

			if response.StatusCode != test.wantStatus {
				t.Errorf("Got status %d, wanted %d", response.StatusCode, test.wantStatus)
			}

			if test.wantError == "" {
				var usersMap map[string][]dto.UserGetDto
				if err := json.NewDecoder(response.Body).Decode(&usersMap); err != nil {
					t.Errorf("Cannot decode body: %v", err)
				}

				users := usersMap["persons"]
				if !slices.Equal(users, test.wantUsers) {
					t.Errorf("Got users %v, wanted %v", users, test.wantUsers)
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

func TestGetCurrentUser(t *testing.T) {
	currentUser := factory.NewTestUser()
	currentAdminUser := factory.NewTestUser(func(u *model.User) { u.Role = model.RoleAdmin })

	tests := []struct {
		name        string
		currentUser model.User
		wantStatus  int
		wantUser    dto.UserGetDto
	}{
		{
			name:        "Returns current user",
			currentUser: currentUser,
			wantStatus:  http.StatusOK,
			wantUser:    currentUser.ToDto(),
		},
		{
			name:        "Returns current admin user",
			currentUser: currentAdminUser,
			wantStatus:  http.StatusOK,
			wantUser:    currentAdminUser.ToDto(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userListSaverMock := &userListSaverMock{}
			preAuthorizerMock := &preAuthorizerMock{}
			router := newTestRouter(t, userListSaverMock, preAuthorizerMock)

			request := httptest.NewRequest(http.MethodGet, "/user", nil)
			request = request.WithContext(context.WithValue(request.Context(), constant.UserContextKey, &test.currentUser))
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)

			response := recorder.Result()
			defer response.Body.Close()

			if response.StatusCode != test.wantStatus {
				t.Errorf("Got status %d, wanted %d", response.StatusCode, test.wantStatus)
			}

			var user dto.UserGetDto
			if err := json.NewDecoder(response.Body).Decode(&user); err != nil {
				t.Errorf("Cannot decode body: %v", err)
			}

			if user != test.currentUser.ToDto() {
				t.Errorf("Got user %v, wanted %v", user, test.wantUser)
			}
		})
	}
}

func TestRegisterUsers(t *testing.T) {
	userPostDto1 := factory.NewTestUserPostDto()
	userPostDto2 := factory.NewTestUserPostDto(func(upd *dto.UserPostDto) {
		upd.Email = "test.teacher@help.com"
		upd.FirstName = "Test"
		upd.LastName = "Teacher"
		upd.Role = "ROLE_TEACHER"
	})
	user1 := model.UserFromPostDto(userPostDto1, "hashed-password")
	user2 := model.UserFromPostDto(userPostDto2, "hashed-password")

	tests := []struct {
		name                  string
		users                 []model.User
		userPostDtos          []dto.UserPostDto
		wantStatus            int
		wantError             string
		wantCreateUserCalls   []model.User
		wantUserByEmailCalls  []string
		wantUserByUuidCalls   []uuid.UUID
		wantHashPasswordCalls []string
		getUserByEmailError   error
		hashPasswordError     error
		createUserError       error
		getUserByUuidError    error
	}{
		{
			name:                  "Registers one user",
			userPostDtos:          []dto.UserPostDto{userPostDto1},
			wantStatus:            http.StatusOK,
			wantCreateUserCalls:   []model.User{user1},
			wantUserByUuidCalls:   []uuid.UUID{user1.Uuid},
			wantUserByEmailCalls:  []string{userPostDto1.Email},
			wantHashPasswordCalls: []string{userPostDto1.Password},
		},
		{
			name:                  "Registers two users",
			userPostDtos:          []dto.UserPostDto{userPostDto1, userPostDto2},
			wantStatus:            http.StatusOK,
			wantCreateUserCalls:   []model.User{user1, user2},
			wantUserByUuidCalls:   []uuid.UUID{user1.Uuid, user2.Uuid},
			wantUserByEmailCalls:  []string{userPostDto1.Email, userPostDto2.Email},
			wantHashPasswordCalls: []string{userPostDto1.Password, userPostDto2.Password},
		},
		{
			name:       "Empty body",
			wantStatus: http.StatusBadRequest,
			wantError:  "Bad request: Cannot decode body as JSON",
		},
		{
			name:                 "One email already in use",
			users:                []model.User{user1},
			userPostDtos:         []dto.UserPostDto{userPostDto1, userPostDto2},
			wantStatus:           http.StatusConflict,
			wantError:            fmt.Sprintf("User with emails '%s' already exists", user1.Email),
			wantUserByEmailCalls: []string{userPostDto1.Email, userPostDto2.Email},
		},
		{
			name:                 "Cannot check if email already in use",
			getUserByEmailError:  fmt.Errorf("Cannot get user by email"),
			userPostDtos:         []dto.UserPostDto{userPostDto1, userPostDto2},
			wantStatus:           http.StatusInternalServerError,
			wantError:            "Internal server error, please try again later",
			wantUserByEmailCalls: []string{userPostDto1.Email},
		},
		{
			name:                  "Cannot hash password",
			hashPasswordError:     fmt.Errorf("Cannot hash password"),
			userPostDtos:          []dto.UserPostDto{userPostDto1, userPostDto2},
			wantStatus:            http.StatusInternalServerError,
			wantError:             "Internal server error, please try again later",
			wantUserByEmailCalls:  []string{userPostDto1.Email, userPostDto2.Email},
			wantHashPasswordCalls: []string{userPostDto1.Password},
		},
		{
			name:                  "Cannot create user",
			createUserError:       fmt.Errorf("Cannot create user"),
			userPostDtos:          []dto.UserPostDto{userPostDto1, userPostDto2},
			wantCreateUserCalls:   []model.User{user1},
			wantStatus:            http.StatusInternalServerError,
			wantError:             "Internal server error, please try again later",
			wantUserByEmailCalls:  []string{userPostDto1.Email, userPostDto2.Email},
			wantHashPasswordCalls: []string{userPostDto1.Password},
		},
		{
			name:                  "Cannot get created user",
			getUserByUuidError:    fmt.Errorf("Cannot get user by UUID"),
			userPostDtos:          []dto.UserPostDto{userPostDto1, userPostDto2},
			wantCreateUserCalls:   []model.User{user1},
			wantStatus:            http.StatusInternalServerError,
			wantError:             "Internal server error, please try again later",
			wantUserByEmailCalls:  []string{userPostDto1.Email, userPostDto2.Email},
			wantUserByUuidCalls:   []uuid.UUID{user1.Uuid},
			wantHashPasswordCalls: []string{userPostDto1.Password},
		},
		{
			name:         "Missing FirstName",
			userPostDtos: []dto.UserPostDto{factory.NewTestUserPostDto(func(upd *dto.UserPostDto) { upd.FirstName = "" })},
			wantStatus:   http.StatusBadRequest,
			wantError:    "Bad request: Key: '[0].FirstName' Error:Field validation for 'FirstName' failed on the 'required' tag",
		},
		{
			name:         "Missing LastName",
			userPostDtos: []dto.UserPostDto{factory.NewTestUserPostDto(func(upd *dto.UserPostDto) { upd.LastName = "" })},
			wantStatus:   http.StatusBadRequest,
			wantError:    "Bad request: Key: '[0].LastName' Error:Field validation for 'LastName' failed on the 'required' tag",
		},
		{
			name:         "Missing Email",
			userPostDtos: []dto.UserPostDto{factory.NewTestUserPostDto(func(upd *dto.UserPostDto) { upd.Email = "" })},
			wantStatus:   http.StatusBadRequest,
			wantError:    "Bad request: Key: '[0].Email' Error:Field validation for 'Email' failed on the 'required' tag",
		},
		{
			name:         "Invalid Email",
			userPostDtos: []dto.UserPostDto{factory.NewTestUserPostDto(func(upd *dto.UserPostDto) { upd.Email = "invalid-email" })},
			wantStatus:   http.StatusBadRequest,
			wantError:    "Bad request: Key: '[0].Email' Error:Field validation for 'Email' failed on the 'email' tag",
		},
		{
			name:         "Missing Password",
			userPostDtos: []dto.UserPostDto{factory.NewTestUserPostDto(func(upd *dto.UserPostDto) { upd.Password = "" })},
			wantStatus:   http.StatusBadRequest,
			wantError:    "Bad request: Key: '[0].Password' Error:Field validation for 'Password' failed on the 'required' tag",
		},
		{
			name:         "Password too simple",
			userPostDtos: []dto.UserPostDto{factory.NewTestUserPostDto(func(upd *dto.UserPostDto) { upd.Password = "short" })},
			wantStatus:   http.StatusBadRequest,
			wantError:    "Bad request: Key: '[0].Password' Error:Field validation for 'Password' failed on the 'password' tag",
		},
		{
			name:         "Password too short",
			userPostDtos: []dto.UserPostDto{factory.NewTestUserPostDto(func(upd *dto.UserPostDto) { upd.Password = "Admin1;" })},
			wantStatus:   http.StatusBadRequest,
			wantError:    "Bad request: Key: '[0].Password' Error:Field validation for 'Password' failed on the 'min' tag",
		},
		{
			name:         "Missing role",
			userPostDtos: []dto.UserPostDto{factory.NewTestUserPostDto(func(upd *dto.UserPostDto) { upd.Role = "" })},
			wantStatus:   http.StatusBadRequest,
			wantError:    "Bad request: Key: '[0].Role' Error:Field validation for 'Role' failed on the 'required' tag",
		},
		{
			name:         "Invalid role",
			userPostDtos: []dto.UserPostDto{factory.NewTestUserPostDto(func(upd *dto.UserPostDto) { upd.Role = "invalid-role" })},
			wantStatus:   http.StatusBadRequest,
			wantError:    "Bad request: Key: '[0].Role' Error:Field validation for 'Role' failed on the 'oneof' tag",
		},
		{
			name:         "Missing group",
			userPostDtos: []dto.UserPostDto{factory.NewTestUserPostDto(func(upd *dto.UserPostDto) { upd.Group = "" })},
			wantStatus:   http.StatusBadRequest,
			wantError:    "Bad request: Key: '[0].Group' Error:Field validation for 'Group' failed on the 'required' tag",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userListSaverMock := &userListSaverMock{users: test.users, getUserByEmailError: test.getUserByEmailError, createUserError: test.createUserError, getUserByUuidError: test.getUserByUuidError}
			preAuthorizerMock := &preAuthorizerMock{}
			router := newTestRouter(t, userListSaverMock, preAuthorizerMock)

			userUuidGeneratorCall := 0
			model.UserUuidGenerator = func() uuid.UUID {
				defer func() { userUuidGeneratorCall++ }()
				return test.wantCreateUserCalls[userUuidGeneratorCall].Uuid
			}
			defer func() { model.UserUuidGenerator = func() uuid.UUID { return uuid.New() } }()

			originalHashPassword := utils.HashPassword
			var hashPasswordCalls []string
			defer func() { utils.HashPassword = originalHashPassword }()
			utils.HashPassword = func(password string) (string, error) {
				hashPasswordCalls = append(hashPasswordCalls, password)
				return "hashed-password", test.hashPasswordError
			}

			var body io.Reader
			if test.userPostDtos != nil {
				userPostDtoJson, err := json.Marshal(test.userPostDtos)
				if err != nil {
					t.Errorf("Cannot marshal userPostDto %v: %v", test.userPostDtos, err)
				}
				body = bytes.NewReader(userPostDtoJson)
			}

			request := httptest.NewRequest(http.MethodPost, "/persons", body)
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)

			response := recorder.Result()
			defer response.Body.Close()

			if response.StatusCode != test.wantStatus {
				t.Errorf("Got status %d, wanted %d", response.StatusCode, test.wantStatus)
			}

			if !slices.Equal(test.wantUserByEmailCalls, userListSaverMock.userByEmailCalls) {
				t.Errorf("got userByEmailCalls %v, wanted %v", userListSaverMock.userByEmailCalls, test.wantUserByEmailCalls)
			}

			if !slices.Equal(test.wantUserByUuidCalls, userListSaverMock.userByUuidCalls) {
				t.Errorf("got userByUuidCalls %v, wanted %v", userListSaverMock.userByUuidCalls, test.wantUserByUuidCalls)
			}

			if !slices.Equal(test.wantCreateUserCalls, userListSaverMock.createUserCalls) {
				t.Errorf("got createUserCalls %v, wanted %v", userListSaverMock.createUserCalls, test.wantCreateUserCalls)
			}

			if !slices.Equal(test.wantHashPasswordCalls, hashPasswordCalls) {
				t.Errorf("got hashPasswordCalls %v, wanted %v", hashPasswordCalls, test.wantHashPasswordCalls)
			}

			if test.wantError == "" {
				var users []dto.UserGetDto
				if err := json.NewDecoder(response.Body).Decode(&users); err != nil {
					t.Errorf("Cannot decode body: %v", err)
				}

				wantUserDtos := dto.ModelsToDtos(test.wantCreateUserCalls)
				if !slices.Equal(users, wantUserDtos) {
					t.Errorf("Got users %v, wanted %v", users, wantUserDtos)
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
