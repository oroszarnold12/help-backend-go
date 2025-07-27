package service

import (
	"context"
	"encoding/json"
	"fmt"
	"help/constant"
	"help/dto"
	"help/model"
	"help/test/factory"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/gorilla/mux"
)

type participationListerMock struct {
	participations []model.Participation
	error          error
	calledUserId   int
	calledCourseId int
}

func (lister *participationListerMock) GetParticipationsOfUser(userId int) ([]model.Participation, error) {
	lister.calledUserId = userId
	return lister.participations, lister.error
}

func (lister *participationListerMock) GetParticipationsOfCourse(courseId int) ([]model.Participation, error) {
	lister.calledCourseId = courseId
	return lister.participations, lister.error
}

func TestGetParticipations(t *testing.T) {
	currentUser := factory.NewTestUser()
	participation1 := factory.NewTestParticipation(func(p *model.Participation) { p.User = currentUser })
	participation2 := factory.NewTestParticipation(func(p *model.Participation) { p.User = currentUser })

	tests := []struct {
		name               string
		participations     []model.Participation
		error              error
		wantStatus         int
		wantParticipations []dto.PariticipationGetDto
		wantError          string
	}{
		{
			name:       "No participations",
			wantStatus: http.StatusOK,
		},
		{
			name:               "Returns participations",
			participations:     []model.Participation{participation1, participation2},
			wantStatus:         http.StatusOK,
			wantParticipations: dto.ModelsToDtos([]model.Participation{participation1, participation2}),
		},
		{
			name:       "Error getting participations",
			error:      fmt.Errorf("Cannot get participations"),
			wantStatus: http.StatusInternalServerError,
			wantError:  "Internal server error, please try again later",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			participationListerMock := &participationListerMock{participations: test.participations, error: test.error}
			participationService := NewParticipationService(participationListerMock)

			router := mux.NewRouter()
			participationService.RegisterRoutes(router)

			request := httptest.NewRequest(http.MethodGet, "/participations", nil)
			request = request.WithContext(context.WithValue(request.Context(), constant.UserContextKey, &currentUser))
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)

			response := recorder.Result()
			defer response.Body.Close()

			if response.StatusCode != test.wantStatus {
				t.Errorf("Got status %d, wanted %d", response.StatusCode, test.wantStatus)
			}

			if test.wantError == "" {
				var participations []dto.PariticipationGetDto
				if err := json.NewDecoder(response.Body).Decode(&participations); err != nil {
					t.Errorf("Cannot decode body: %v", err)
				}

				if !slices.Equal(participations, test.wantParticipations) {
					t.Errorf("Got participations %v, wanted %v", participations, test.wantParticipations)
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

			if participationListerMock.calledUserId != currentUser.Id {
				t.Errorf("GetParticipationsOfUser got called with userId %d, wanted %d", participationListerMock.calledUserId, currentUser.Id)
			}
		})
	}
}

func TestGetParticipants(t *testing.T) {
	currentUser := factory.NewTestUser()
	user1 := factory.NewTestUser(func(u *model.User) { u.Id = 2 })
	user2 := factory.NewTestUser(func(u *model.User) { u.Id = 3 })
	course := factory.NewTestCourse()
	participation1 := factory.NewTestParticipation(func(p *model.Participation) {
		p.User = currentUser
		p.Course = course
	})
	participation2 := factory.NewTestParticipation(func(p *model.Participation) {
		p.User = user1
		p.Course = course
	})
	participation3 := factory.NewTestParticipation(func(p *model.Participation) {
		p.User = user2
		p.Course = course
	})

	tests := []struct {
		name               string
		path               string
		participations     []model.Participation
		error              error
		wantStatus         int
		wantParticipants   []dto.UserGetDto
		wantError          string
		wantCalledCourseId int
	}{
		{
			name:               "Returns user as participant",
			path:               fmt.Sprintf("/courses/%d/participants", course.Id),
			participations:     []model.Participation{participation1},
			wantStatus:         http.StatusOK,
			wantParticipants:   []dto.UserGetDto{currentUser.ToDto()},
			wantCalledCourseId: course.Id,
		},
		{
			name:               "Returns all users as participant",
			path:               fmt.Sprintf("/courses/%d/participants", course.Id),
			participations:     []model.Participation{participation1, participation2, participation3},
			wantStatus:         http.StatusOK,
			wantParticipants:   []dto.UserGetDto{currentUser.ToDto(), user1.ToDto(), user2.ToDto()},
			wantCalledCourseId: course.Id,
		},
		{
			name:           "Invalid course ID returns bad request error",
			path:           "/courses/invalid6/participants",
			participations: []model.Participation{participation1},
			wantStatus:     http.StatusBadRequest,
			wantError:      "Bad request: Invalid course ID",
		},
		{
			name:               "Error getting participations",
			path:               fmt.Sprintf("/courses/%d/participants", course.Id),
			error:              fmt.Errorf("Cannot get participations"),
			wantStatus:         http.StatusInternalServerError,
			wantError:          "Internal server error, please try again later",
			wantCalledCourseId: course.Id,
		},
		{
			name:               "Returns 403 if the user is not participating on the course",
			path:               fmt.Sprintf("/courses/%d/participants", course.Id),
			participations:     []model.Participation{participation2, participation3},
			wantStatus:         http.StatusForbidden,
			wantError:          "You do not have the required role to perform this action",
			wantCalledCourseId: course.Id,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			participationListerMock := &participationListerMock{participations: test.participations, error: test.error}
			participationService := NewParticipationService(participationListerMock)

			router := mux.NewRouter()
			participationService.RegisterRoutes(router)

			request := httptest.NewRequest(http.MethodGet, test.path, nil)
			request = request.WithContext(context.WithValue(request.Context(), constant.UserContextKey, &currentUser))
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)

			response := recorder.Result()
			defer response.Body.Close()

			if response.StatusCode != test.wantStatus {
				t.Errorf("Got status %d, wanted %d", response.StatusCode, test.wantStatus)
			}

			if test.wantError == "" {
				var participants []dto.UserGetDto
				if err := json.NewDecoder(response.Body).Decode(&participants); err != nil {
					t.Errorf("Cannot decode body: %v", err)
				}

				if !slices.Equal(participants, test.wantParticipants) {
					t.Errorf("Got participants %v, wanted %v", participants, test.wantParticipants)
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

			if test.wantCalledCourseId != 0 && participationListerMock.calledCourseId != test.wantCalledCourseId {
				t.Errorf("GetParticipationsOfCourse got called with courseId %d, wanted %d", participationListerMock.calledCourseId, course.Id)
			}
		})
	}
}
