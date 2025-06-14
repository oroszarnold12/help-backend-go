package dto

type PariticipationGetDto struct {
	Course          ThinCourseGetDto `json:"course"`
	ShowOnDashboard bool             `json:"showOnDashboard"`
}
