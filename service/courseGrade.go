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

type CourseGradeService struct {
	assignmentGradeDao *dao.AssignmentGradeDao
}

func NewCourseGradeService(assignmentGradeDao *dao.AssignmentGradeDao) *CourseGradeService {
	return &CourseGradeService{assignmentGradeDao: assignmentGradeDao}
}

func (service *CourseGradeService) RegisterRoutes(authorizedRouter *mux.Router) {
	authorizedRouter.HandleFunc("/courses/{id}/grades", service.getGrades).Methods(http.MethodGet)
}

func (service *CourseGradeService) getGrades(writer http.ResponseWriter, request *http.Request) {
	user := utils.GetCurrentUser(request)

	courseIdString := mux.Vars(request)["id"]
	courseId, err := strconv.Atoi(courseIdString)
	if err != nil {
		utils.WriteError(writer, fmt.Errorf("%w: %v", errorsx.NewBadRequestError("Invalid course ID"), err))
		return
	}

	assignmentGrades, err := service.assignmentGradeDao.GetAssignmentGrades(courseId, user.Id)
	if err != nil {
		utils.WriteError(writer, err)
	}

	utils.WriteJson(writer, http.StatusOK, dto.CourseGradeGetDto{AssignmentGrades: dto.ModelsToDtos(assignmentGrades)})
}
