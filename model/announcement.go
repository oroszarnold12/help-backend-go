package model

import (
	"help/dto"
	"time"

	"github.com/google/uuid"
)

type Announcement struct {
	Id      int
	Uuid    uuid.UUID
	Name    string
	Date    time.Time
	Content string
	Course  Course
	Creator User
}

func (announcement Announcement) ToThinDto() dto.ThinAnnouncementGetDto {
	return dto.ThinAnnouncementGetDto{
		Id:      announcement.Id,
		Uuid:    announcement.Uuid.String(),
		Name:    announcement.Name,
		Content: announcement.Content,
		Date:    announcement.Date.String(),
	}
}
