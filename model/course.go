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
	Announcements []Announcement
}

func (course Course) ToDto() dto.CourseGetDto {
	return dto.CourseGetDto{
		Id:            course.Id,
		Uuid:          course.Uuid.String(),
		Name:          course.Name,
		LongName:      course.LongName,
		Descirption:   course.Description,
		Teacher:       course.Teacher.ToDto(),
		Announcements: dto.ModelsToThinDtos(course.Announcements),
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
