package dto

type CourseGetDto struct {
	Id            int                      `json:"id"`
	Uuid          string                   `json:"uuid"`
	Name          string                   `json:"name"`
	LongName      string                   `json:"longName"`
	Descirption   string                   `json:"description"`
	Teacher       UserGetDto               `json:"teacher"`
	Assignments   []ThinAssignmentGetDto   `json:"assignments"`
	Announcements []ThinAnnouncementGetDto `json:"announcements"`
	Discussions   []ThinDiscussionGetDto   `json:"discussions"`
	Quizzes       []ThingQuizGetDto        `json:"quizzes"`
}

type ThinCourseGetDto struct {
	Id       int    `json:"id"`
	Uuid     string `json:"uuid"`
	Name     string `json:"name"`
	LongName string `json:"longName"`
}

type CourseGradeGetDto struct {
	AssignmentGrades []AssignemntGradeGetDto `json:"assignmentGrades"`
}
