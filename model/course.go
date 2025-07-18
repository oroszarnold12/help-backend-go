package model

import (
	"github.com/google/uuid"
	"help/dto"
)

type Course struct {
	Id            int
	Uuid          uuid.UUID
	Name          string
	LongName      string
	Description   string
	Teacher       User
	Assignments   []Assignment
	Announcements []Announcement
	Discussions   []Discussion
}

func (course Course) ToDto() dto.CourseGetDto {
	return dto.CourseGetDto{
		Id:            course.Id,
		Uuid:          course.Uuid.String(),
		Name:          course.Name,
		LongName:      course.LongName,
		Descirption:   course.Description,
		Teacher:       course.Teacher.ToDto(),
		Assignments:   dto.ModelsToThinDtos(course.Assignments),
		Announcements: dto.ModelsToThinDtos(course.Announcements),
		Discussions:   dto.ModelsToThinDtos(course.Discussions),
	}
}

func (course Course) ToThinDto() dto.ThinCourseGetDto {
	return dto.ThinCourseGetDto{
		Id:       course.Id,
		Uuid:     course.Uuid.String(),
		Name:     course.Name,
		LongName: course.LongName,
	}
}
