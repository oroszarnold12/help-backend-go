package dao

import (
	"database/sql"
	"fmt"
	"help/errorsx"
	"help/model"
	"strconv"
)

const courseSelectFields = `
	c.id, c.uuid, c.name, c.long_name, c.description,
	u.id, u.uuid, u.first_name, u.last_name, u.email, u.role, u.password
`

const announcementFields = `
	a.id, a.uuid, a.name, a.content, a.date
`

const assignmentFields = `
	a.id, a.uuid, a.name, a.due_date, a.points, a.published
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

	courses, err = completeCourseModel(dao, courses)
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

	courses, err = completeCourseModel(dao, courses)
	if err != nil {
		return nil, err
	}

	if len(courses) == 0 {
		return nil, errorsx.NewNotFoundError("Course", strconv.Itoa(courseId))
	}

	return &courses[0], nil
}

func (dao *CourseDao) getAssignmentsOfCourse(courseId int) ([]model.Assignment, error) {
	rows, err := dao.db.Query(
		fmt.Sprintf(`
			SELECT %s
			FROM assignments a
			WHERE a.course_id = ?
			`,
			assignmentFields,
		),
		courseId,
	)

	if err != nil {
		return nil, fmt.Errorf("Cannot query db: %w", err)
	}
	defer rows.Close()

	assignemnts, err := scanRowstoAssignments(rows)
	if err != nil {
		return nil, err
	}

	return assignemnts, nil
}

func (dao *CourseDao) getAnnouncementsOfCourse(courseId int) ([]model.Announcement, error) {
	rows, err := dao.db.Query(
		fmt.Sprintf(`
			SELECT %s
			FROM announcements a
			WHERE a.course_id = ?
			`,
			announcementFields,
		),
		courseId,
	)

	if err != nil {
		return nil, fmt.Errorf("Cannot query db: %w", err)
	}
	defer rows.Close()

	announcements, err := scanRowstoAnnouncement(rows)
	if err != nil {
		return nil, err
	}

	return announcements, nil
}

func completeCourseModel(dao *CourseDao, courses []model.Course) ([]model.Course, error) {
	for index := range courses {
		courseId := courses[index].Id
		assignments, err := dao.getAssignmentsOfCourse(courseId)
		if err != nil {
			return nil, fmt.Errorf("Cannot get assignments of course '%d': %w", courseId, err)
		}

		announcements, err := dao.getAnnouncementsOfCourse(courseId)
		if err != nil {
			return nil, fmt.Errorf("Cannot get announcements of course '%d': %w", courseId, err)
		}

		courses[index].Assignments = assignments
		courses[index].Announcements = announcements
	}

	return courses, nil
}

func scanRowstoAssignments(rows *sql.Rows) ([]model.Assignment, error) {
	var assignments []model.Assignment

	for rows.Next() {
		var assignment model.Assignment
		err := rows.Scan(
			&assignment.Id,
			&assignment.Uuid,
			&assignment.Name,
			&assignment.DueDate,
			&assignment.Points,
			&assignment.Published,
		)

		if err != nil {
			return nil, fmt.Errorf("Cannot scan assignment rows into model: %w", err)
		}

		assignments = append(assignments, assignment)
	}

	return assignments, nil
}

func scanRowstoAnnouncement(rows *sql.Rows) ([]model.Announcement, error) {
	var announcements []model.Announcement

	for rows.Next() {
		var announcement model.Announcement
		err := rows.Scan(
			&announcement.Id,
			&announcement.Uuid,
			&announcement.Name,
			&announcement.Content,
			&announcement.Date,
		)

		if err != nil {
			return nil, fmt.Errorf("Cannot scan announcement rows into model: %w", err)
		}

		announcements = append(announcements, announcement)
	}

	return announcements, nil
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

		courses = append(courses, course)
	}

	return courses, nil
}
