package service

import (
	"help/dao"
	"help/utils"
	"net/http"

	"github.com/gorilla/mux"
)

type GroupService struct {
	userLister dao.UserLister
}

func NewGroupService(userLister dao.UserLister) *GroupService {
	return &GroupService{userLister: userLister}
}

func (service *GroupService) RegisterRoutes(authorizedRouter *mux.Router) {
	authorizedRouter.HandleFunc("/groups", service.getGroups).Methods(http.MethodGet)
}

func (service *GroupService) getGroups(writer http.ResponseWriter, request *http.Request) {
	users, err := service.userLister.GetUsers()
	if err != nil {
		utils.WriteError(writer, err)
	}

	groupSet := map[string]struct{}{}
	for _, user := range users {
		if user.Group != "" {
			groupSet[user.Group] = struct{}{}
		}
	}

	groups := make([]string, 0, len(groupSet))
	for group := range groupSet {
		groups = append(groups, group)
	}

	utils.WriteJson(writer, http.StatusOK, map[string][]string{"personGroups": groups})
}
