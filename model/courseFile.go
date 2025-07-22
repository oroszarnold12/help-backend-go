package model

import (
	"help/dto"
	"time"

	"github.com/google/uuid"
)

type CourseFile struct {
	Id           int
	Uuid         uuid.UUID
	Name         string
	Size         int
	CreationDate time.Time
	Uploader     User
}

func (courseFile CourseFile) ToDto() dto.CourseFileGetDto {
	return dto.CourseFileGetDto{
		Id:           courseFile.Id,
		Uuid:         courseFile.Uuid.String(),
		Name:         courseFile.Name,
		Size:         courseFile.Size,
		CreationDate: courseFile.CreationDate.String(),
		Uploader:     courseFile.Uploader.ToDto(),
	}
}
