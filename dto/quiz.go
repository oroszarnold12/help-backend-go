package dto

type ThingQuizGetDto struct {
	Id        int     `json:"id"`
	Uuid      string  `json:"uuid"`
	Name      string  `json:"name"`
	DueDate   string  `json:"dueDate"`
	Points    float64 `json:"points"`
	Published bool    `json:"published"`
}
