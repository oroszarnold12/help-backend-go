package dto

type ThinAnnouncementGetDto struct {
	Id      int    `json:"id"`
	Uuid    string `json:"uuid"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Date    string `json:"date"`
}
