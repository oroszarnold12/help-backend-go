package service

import (
	"fmt"
	"help/dao"
	"help/dto"
	"help/errorsx"
	"help/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type InvitationService struct {
	invitationDao    *dao.InvitationDao
	participationDao *dao.ParticipationDao
}

func NewInvitationService(invitationDao *dao.InvitationDao, participationDao *dao.ParticipationDao) *InvitationService {
	return &InvitationService{invitationDao: invitationDao, participationDao: participationDao}
}

func (service *InvitationService) RegisterRoutes(authorizedRouter *mux.Router) {
	authorizedRouter.HandleFunc("/invitations", service.getInvitations).Methods(http.MethodGet)
	authorizedRouter.HandleFunc("/invitations/{id}", service.deleteInvitation).Methods(http.MethodDelete)
}

func (service *InvitationService) getInvitations(writer http.ResponseWriter, request *http.Request) {
	user := utils.GetCurrentUser(request)

	invitations, err := service.invitationDao.GetInvitationsOfUser(user.Id)
	if err != nil {
		utils.WriteError(writer, err)
		return
	}

	utils.WriteJson(writer, http.StatusOK, map[string][]dto.InvitationGetDto{"invitations": dto.ModelsToDtos(invitations)})
}

func (service *InvitationService) deleteInvitation(writer http.ResponseWriter, request *http.Request) {
	user := utils.GetCurrentUser(request)

	invitationIdString := mux.Vars(request)["id"]
	invitationId, err := strconv.Atoi(invitationIdString)
	if err != nil {
		utils.WriteError(writer, fmt.Errorf("%w: %v", errorsx.NewBadRequestError("Invalid invitation ID"), err))
		return
	}

	invitation, err := service.invitationDao.GetInvitationOfUser(user.Id, invitationId)
	if err != nil {
		utils.WriteError(writer, err)
		return
	}

	accept := request.URL.Query().Get("accept") == "true"

	if accept {
		err := service.participationDao.CreateParticipation(user.Id, invitation.Course.Id)
		if err != nil {
			utils.WriteError(writer, fmt.Errorf("Cannot create participation for user '%v' on course '%v': %w", user.Id, invitation.Course.Id, err))
			return
		}
	}

	err = service.invitationDao.DeleteInvitation(invitationId)
	if err != nil {
		utils.WriteError(writer, err)
		return
	}

	utils.WriteJson(writer, http.StatusOK, nil)
}
