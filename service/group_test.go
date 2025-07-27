package service

import (
	"encoding/json"
	"fmt"
	"help/model"
	"help/test/factory"

	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/gorilla/mux"
)

type userListerMock struct {
	users []model.User
	error error
}

func (lister *userListerMock) GetUsers() ([]model.User, error) {
	return lister.users, lister.error
}

func TestGetGroups(t *testing.T) {
	user1 := factory.NewTestUser(func(u *model.User) { u.Group = "Group1" })
	user2 := factory.NewTestUser(func(u *model.User) { u.Group = "Group2" })
	user3 := factory.NewTestUser(func(u *model.User) { u.Group = "Group2" })
	user4 := factory.NewTestUser(func(u *model.User) { u.Group = "" })
	user5 := factory.NewTestUser(func(u *model.User) { u.Group = "" })

	tests := []struct {
		name       string
		users      []model.User
		error      error
		wantGroups []string
		wantStatus int
		wantError  string
	}{
		{
			name:       "No users",
			users:      []model.User{},
			wantGroups: []string{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "No users with groups",
			users:      []model.User{user4, user5},
			wantGroups: []string{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "Users with groups",
			users:      []model.User{user1, user2, user3, user4, user5},
			wantGroups: []string{"Group1", "Group2"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "Error getting users",
			error:      fmt.Errorf("Cannot get users"),
			wantStatus: http.StatusInternalServerError,
			wantError:  "Internal server error, please try again later",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userListerMock := &userListerMock{users: test.users, error: test.error}
			groupService := NewGroupService(userListerMock)

			router := mux.NewRouter()
			groupService.RegisterRoutes(router)

			request := httptest.NewRequest(http.MethodGet, "/groups", nil)
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)

			response := recorder.Result()
			defer response.Body.Close()

			if response.StatusCode != test.wantStatus {
				t.Errorf("Got status %d, wanted %d", response.StatusCode, test.wantStatus)
			}

			if test.wantError == "" {
				var data map[string][]string
				if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
					t.Errorf("Cannot decode body: %v", err)
				}

				groups := data["personGroups"]
				slices.Sort(groups)
				if !slices.Equal(groups, test.wantGroups) {
					t.Errorf("Got groups %v, wanted %v", groups, test.wantGroups)
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
