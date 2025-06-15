package service

import (
	"help/dao"
	"help/dto"
	"help/utils"
	"net/http"

	"github.com/gorilla/mux"
)

type CourseSerivce struct {
	courseDao *dao.CourseDao
}

func NewCourseService(courseDao *dao.CourseDao) *CourseSerivce {
	return &CourseSerivce{courseDao: courseDao}
}

func (service *CourseSerivce) RegisterRoutes(authorizedRouter *mux.Router) {
	authorizedRouter.HandleFunc("/courses", service.getCourses).Methods("GET")
}

func (service *CourseSerivce) getCourses(writer http.ResponseWriter, request *http.Request) {
	user := utils.GetCurrentUser(request)

	courses, err := service.courseDao.GetCoursesOfUser(user.Id)
	if err != nil {
		utils.WriteError(writer, err)
		return
	}

	utils.WriteJson(writer, http.StatusOK, map[string][]dto.ThinCourseGetDto{"courses": dto.ModelsToThinDtos(courses)})
}
