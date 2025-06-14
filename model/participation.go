package model

import (
	"help/dto"

	"github.com/google/uuid"
)

type Participation struct {
	Id              int
	Uuid            uuid.UUID
	ShowOnDashboard bool
	Course          Course
	User            User
}

func (participation Participation) ToDto() dto.PariticipationGetDto {
	return dto.PariticipationGetDto{
		Course:          participation.Course.ToThinDto(),
		ShowOnDashboard: participation.ShowOnDashboard,
	}
}
