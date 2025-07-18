package model

import (
	"help/dto"
	"time"

	"github.com/google/uuid"
)

type Quiz struct {
	Id        int
	Uuid      uuid.UUID
	Name      string
	DueDate   time.Time
	Points    float64
	Published bool
}

func (quiz Quiz) ToThinDto() dto.ThingQuizGetDto {
	return dto.ThingQuizGetDto{
		Id:        quiz.Id,
		Uuid:      quiz.Uuid.String(),
		Name:      quiz.Name,
		DueDate:   quiz.DueDate.String(),
		Points:    quiz.Points,
		Published: quiz.Published,
	}
}
