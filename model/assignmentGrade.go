package model

import (
	"help/dto"

	"github.com/google/uuid"
)

type AssignmentGrade struct {
	Id         int
	Uuid       uuid.UUID
	Submitter  User
	Assignment Assignment
	Grade      float64
}

func (grade AssignmentGrade) ToDto() dto.AssignemntGradeGetDto {
	return dto.AssignemntGradeGetDto{
		Id:         grade.Id,
		Uuid:       grade.Uuid.String(),
		Submitter:  grade.Submitter.ToDto(),
		Assignment: grade.Assignment.ToThinDto(),
		Grade:      grade.Grade,
	}
}
