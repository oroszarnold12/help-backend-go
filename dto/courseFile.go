package dto

type CourseFileGetDto struct {
	Id           int        `json:"id"`
	Uuid         string     `json:"uuid"`
	Name         string     `json:"fileName"`
	Size         int        `json:"size"`
	CreationDate string     `json:"creationDate"`
	Uploader     UserGetDto `json:"uploader"`
}
