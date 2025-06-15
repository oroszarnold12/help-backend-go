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

type CourseSerivce struct {
	courseDao *dao.CourseDao
}

func NewCourseService(courseDao *dao.CourseDao) *CourseSerivce {
	return &CourseSerivce{courseDao: courseDao}
}

func (service *CourseSerivce) RegisterRoutes(authorizedRouter *mux.Router) {
	authorizedRouter.HandleFunc("/courses", service.getCourses).Methods("GET")
	authorizedRouter.HandleFunc("/courses/{id}", service.getCourse).Methods("GET")
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

func (service *CourseSerivce) getCourse(writer http.ResponseWriter, request *http.Request) {
	user := utils.GetCurrentUser(request)

	courseIdString := mux.Vars(request)["id"]
	courseId, err := strconv.Atoi(courseIdString)
	if err != nil {
		utils.WriteError(writer, fmt.Errorf("%w: %v", errorsx.NewBadRequestError("Invalid course ID"), err))
		return
	}

	course, err := service.courseDao.GetCourseOfUser(courseId, user.Id)
	if err != nil {
		utils.WriteError(writer, err)
		return
	}

	utils.WriteJson(writer, http.StatusOK, course.ToDto())
}
