package service

import (
	"help/dao"
	"help/dto"
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
	authorizedRouter.HandleFunc("/participations", service.getParticipations).Methods(http.MethodGet)
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
