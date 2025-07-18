package dto

type AssignemntGradeGetDto struct {
	Id         int                  `json:"id"`
	Uuid       string               `json:"uuid"`
	Submitter  UserGetDto           `json:"submitter"`
	Assignment ThinAssignmentGetDto `json:"assignment"`
	Grade      float64              `json:"grade"`
}
