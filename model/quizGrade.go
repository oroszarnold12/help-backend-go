package model

import (
	"help/dto"

	"github.com/google/uuid"
)

type QuizGrade struct {
	Id        int
	Uuid      uuid.UUID
	Submitter User
	Quiz      Quiz
	Grade     float64
}

func (grade QuizGrade) ToDto() dto.QuizGradeGetDto {
	return dto.QuizGradeGetDto{
		Id:        grade.Id,
		Uuid:      grade.Uuid.String(),
		Submitter: grade.Submitter.ToDto(),
		Quiz:      grade.Quiz.ToThinDto(),
		Grade:     grade.Grade,
	}
}
