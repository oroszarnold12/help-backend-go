package dto

type InvitationGetDto struct {
	Id     int              `json:"id"`
	Uuid   string           `json:"uuid"`
	Course ThinCourseGetDto `json:"course"`
}
