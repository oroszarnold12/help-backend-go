package dao

import (
	"database/sql"
	"fmt"
	"help/model"
)

const courseSelectFields = `
	c.id, c.uuid, c.name, c.long_name, c.description,
	u.id, u.uuid, u.first_name, u.last_name, u.email, u.role, u.password
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

func scanRowToCourse(row *sql.Row) (*model.Course, error) {
	var course model.Course
	var teacher model.User

	err := row.Scan(
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
	)

	if err != nil {
		return nil, fmt.Errorf("Cannot scan course row into model: %w", err)
	}

	course.Teacher = teacher

	return &course, nil
}

func scanRowsToCourses(rows *sql.Rows) ([]model.Course, error) {
	var courses []model.Course

	for rows.Next() {
		var course model.Course
		var teacher model.User

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
		)

		if err != nil {
			return nil, fmt.Errorf("Cannot scan course row into model: %w", err)
		}

		course.Teacher = teacher

		courses = append(courses, course)
	}

	return courses, nil
}
