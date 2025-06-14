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
	user := request.Context().Value(constant.UserContextKey).(*model.User)

	courses, err := service.courseDao.GetCoursesOfUser(user.Id)
	if err != nil {
		utils.WriteError(writer, err)
		return
	}

	courseDtos := make([]dto.ThinCourseGetDto, len(courses))
	for index := range courses {
		courseDtos[index] = courses[index].ToThinDto()
	}

	utils.WriteJson(writer, http.StatusOK, map[string][]dto.ThinCourseGetDto{"courses": courseDtos})
}
