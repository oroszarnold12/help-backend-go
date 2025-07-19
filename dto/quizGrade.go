package dto

type QuizGradeGetDto struct {
	Id        int             `json:"id"`
	Uuid      string          `json:"uuid"`
	Submitter UserGetDto      `json:"submitter"`
	Quiz      ThingQuizGetDto `json:"quiz"`
	Grade     float64         `json:"grade"`
}
