package model

import (
	"help/dto"
	"time"

	"github.com/google/uuid"
)

type Discussion struct {
	Id      int
	Uuid    uuid.UUID
	Name    string
	Date    time.Time
	Creator User
}

func (discussion Discussion) ToThinDto() dto.ThinDiscussionGetDto {
	return dto.ThinDiscussionGetDto{
		Id:      discussion.Id,
		Uuid:    discussion.Uuid.String(),
		Name:    discussion.Name,
		Date:    discussion.Date.String(),
		Creator: discussion.Creator.ToDto(),
	}
}
