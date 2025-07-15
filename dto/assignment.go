package dto

type ThinAssignmentGetDto struct {
	Id        int    `json:"id"`
	Uuid      string `json:"uuid"`
	Name      string `json:"name"`
	DueDate   string `json:"dueDate"`
	Points    int    `json:"points"`
	Published bool   `json:"published"`
}
