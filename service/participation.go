package service

import (
	"fmt"
	"help/dao"
	"help/dto"
	"help/errorsx"
	"help/model"
	"help/utils"
	"net/http"
	"slices"
	"strconv"

	"github.com/gorilla/mux"
)

type ParticipationService struct {
	participationDao *dao.ParticipationDao
}

func NewParticipaionService(participationDao *dao.ParticipationDao) *ParticipationService {
	return &ParticipationService{participationDao: participationDao}
}

func (service *ParticipationService) RegisterRoutes(authorizedRouter *mux.Router) {
	authorizedRouter.HandleFunc("/participations", service.getParticipations).Methods(http.MethodGet)
	authorizedRouter.HandleFunc("/courses/{id}/participants", service.getParticipants).Methods(http.MethodGet)
}

func (service *ParticipationService) getParticipations(writer http.ResponseWriter, request *http.Request) {
	user := utils.GetCurrentUser(request)

	participations, err := service.participationDao.GetParticipationsOfUser(user.Id)
	if err != nil {
		utils.WriteError(writer, err)
		return
	}

	utils.WriteJson(writer, http.StatusOK, dto.ModelsToDtos(participations))
}

func (service *ParticipationService) getParticipants(writer http.ResponseWriter, request *http.Request) {
	user := utils.GetCurrentUser(request)

	courseIdString := mux.Vars(request)["id"]
	courseId, err := strconv.Atoi(courseIdString)
	if err != nil {
		utils.WriteError(writer, fmt.Errorf("%w: %v", errorsx.NewBadRequestError("Invalid course ID"), err))
		return
	}

	participations, err := service.participationDao.GetParticipationsOfCourse(courseId)
	if err != nil {
		utils.WriteError(writer, err)
		return
	}

	userParticipates := slices.ContainsFunc(participations, func(participation model.Participation) bool {
		return participation.User.Id == user.Id
	})
	if !userParticipates {
		utils.WriteError(writer, errorsx.NewForbiddenError())
		return
	}

	participants := make([]model.User, len(participations))
	for index := range participations {
		participants[index] = participations[index].User
	}

	utils.WriteJson(writer, http.StatusOK, dto.ModelsToDtos(participants))
}
