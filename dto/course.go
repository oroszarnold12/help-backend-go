package dto

type CourseGetDto struct {
	Id            int                      `json:"id"`
	Uuid          string                   `json:"uuid"`
	Name          string                   `json:"name"`
	LongName      string                   `json:"longName"`
	Descirption   string                   `json:"description"`
	Teacher       UserGetDto               `json:"teacher"`
	Announcements []ThinAnnouncementGetDto `json:"announcements"`
}

type ThinCourseGetDto struct {
	Id       int    `json:"id"`
	Uuid     string `json:"uuid"`
	Name     string `json:"name"`
	LongName string `json:"longName"`
}
