package dao

import (
	"database/sql"
	"fmt"
	"help/errorsx"
	"help/model"
	"strconv"

	"github.com/google/uuid"
)

const courseSelectFields = `
	c.id, c.uuid, c.name, c.long_name, c.description,
	u.id, u.uuid, u.first_name, u.last_name, u.email, u.role, u.password,
	a.id, a.uuid, a.name, a.content, a.date
`

type CourseDao struct {
	db *sql.DB
}

func NewCourseDao(db *sql.DB) *CourseDao {
	return &CourseDao{db: db}
}

func (dao *CourseDao) GetCoursesOfUser(userId int) ([]model.Course, error) {
	rows, err := dao.db.Query(
		fmt.Sprintf(`
			SELECT %s 
			FROM participations p 
			JOIN courses c ON c.id = p.course_id 
			JOIN users u on c.teacher_id = u.id 
			LEFT JOIN announcements a on c.id = a.course_id
			WHERE p.user_id = ?`,
			courseSelectFields,
		),
		userId,
	)

	if err != nil {
		return nil, fmt.Errorf("Cannot query db: %w", err)
	}
	defer rows.Close()

	courses, err := scanRowsToCourses(rows)
	if err != nil {
		return nil, err
	}

	return courses, nil
}

func (dao *CourseDao) GetCourseOfUser(courseId int, userId int) (*model.Course, error) {
	rows, err := dao.db.Query(
		fmt.Sprintf(`
			SELECT %s 
			FROM participations p 
			JOIN courses c ON c.id = p.course_id 
			JOIN users u on c.teacher_id = u.id 
			LEFT JOIN announcements a on c.id = a.course_id
			WHERE p.user_id = ? AND c.id = ?`,
			courseSelectFields,
		),
		userId,
		courseId,
	)

	if err != nil {
		return nil, fmt.Errorf("Cannot query db: %w", err)
	}
	defer rows.Close()

	courses, err := scanRowsToCourses(rows)
	if err != nil {
		return nil, err
	}

	if len(courses) == 0 {
		return nil, errorsx.NewNotFoundError("Course", strconv.Itoa(courseId))
	}

	return &courses[0], nil
}

func scanRowsToCourses(rows *sql.Rows) ([]model.Course, error) {
	courseMap := make(map[int]*model.Course)

	for rows.Next() {
		var course model.Course
		var teacher model.User

		var announcementId sql.NullInt64
		var announcementUuid, announcementName, announcementContent sql.NullString
		var announcementDate sql.NullTime

		err := rows.Scan(
			&course.Id,
			&course.Uuid,
			&course.Name,
			&course.LongName,
			&course.Description,
			&teacher.Id,
			&teacher.Uuid,
			&teacher.FirstName,
			&teacher.LastName,
			&teacher.Email,
			&teacher.Role,
			&teacher.Password,
			&announcementId,
			&announcementUuid,
			&announcementName,
			&announcementContent,
			&announcementDate,
		)

		if err != nil {
			return nil, fmt.Errorf("Cannot scan course row into model: %w", err)
		}

		existingCourse, ok := courseMap[course.Id]

		if !ok {
			course.Teacher = teacher
			course.Announcements = []model.Announcement{}
			courseMap[course.Id] = &course
			existingCourse = &course
		}

		if announcementId.Valid {
			announcement := model.Announcement{
				Id:      int(announcementId.Int64),
				Uuid:    uuid.MustParse(announcementUuid.String),
				Name:    announcementName.String,
				Date:    announcementDate.Time,
				Content: announcementContent.String,
			}

			existingCourse.Announcements = append(existingCourse.Announcements, announcement)
		}
	}

	courses := []model.Course{}
	for _, course := range courseMap {
		courses = append(courses, *course)
	}

	return courses, nil
}
