package factory

import (
	"help/model"

	"github.com/google/uuid"
)

func NewTestCourse(overrides ...func(*model.Course)) model.Course {
	course := model.Course{
		Id:          1,
		Uuid:        uuid.New(),
		Name:        "Course",
		LongName:    "Long Course",
		Description: "Description",
		Teacher:     NewTestUser(func(u *model.User) { u.Role = model.RoleTeacher }),
	}

	for _, override := range overrides {
		override(&course)
	}

	return course
}
