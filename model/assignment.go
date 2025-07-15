package model

import (
	"help/dto"
	"time"

	"github.com/google/uuid"
)

type Assignment struct {
	Id        int
	Uuid      uuid.UUID
	Name      string
	DueDate   time.Time
	Points    int
	Published bool
}

func (assignment Assignment) ToThinDto() dto.ThinAssignmentGetDto {
	return dto.ThinAssignmentGetDto{
		Id:        assignment.Id,
		Uuid:      assignment.Uuid.String(),
		Name:      assignment.Name,
		DueDate:   assignment.DueDate.String(),
		Points:    assignment.Points,
		Published: assignment.Published,
	}
}
