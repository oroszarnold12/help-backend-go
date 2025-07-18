package dto

type ThinDiscussionGetDto struct {
	Id      int        `json:"id"`
	Uuid    string     `json:"uuid"`
	Name    string     `json:"name"`
	Date    string     `json:"date"`
	Creator UserGetDto `json:"creator"`
}
