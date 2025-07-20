package model

import (
	"help/dto"

	"github.com/google/uuid"
)

type Invitation struct {
	Id     int
	Uuid   uuid.UUID
	Course Course
	User   User
}

func (invitation Invitation) ToDto() dto.InvitationGetDto {
	return dto.InvitationGetDto{
		Id:     invitation.Id,
		Uuid:   invitation.Uuid.String(),
		Course: invitation.Course.ToThinDto(),
	}
}
