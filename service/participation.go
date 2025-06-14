package service

import (
	"help/constant"
	"help/dao"
	"help/dto"
	"help/model"
	"help/utils"
	"net/http"

	"github.com/gorilla/mux"
)

type ParticipationService struct {
	participationDao *dao.ParticipationDao
}

func NewParticipaionService(participationDao *dao.ParticipationDao) *ParticipationService {
	return &ParticipationService{participationDao: participationDao}
}

func (service *ParticipationService) RegisterRoutes(authorizedRouter *mux.Router) {
	authorizedRouter.HandleFunc("/participations", service.getParticipations).Methods("GET")
}

func (service *ParticipationService) getParticipations(writer http.ResponseWriter, request *http.Request) {
	user := request.Context().Value(constant.UserContextKey).(*model.User)

	participations, err := service.participationDao.GetParticipationsOfUser(user.Id)
	if err != nil {
		utils.WriteError(writer, err)
		return
	}

	participationDtos := make([]dto.PariticipationGetDto, len(participations))
	for index := range participations {
		participationDtos[index] = participations[index].ToDto()
	}

	utils.WriteJson(writer, http.StatusOK, participationDtos)
}
