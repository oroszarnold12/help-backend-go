package dao

import (
	"database/sql"
	"fmt"
	"help/errorsx"
	"help/model"
	"help/utils"
	"strconv"
)

const courseSelectFields = `
	c.id, c.uuid, c.name, c.long_name, c.description,
	u.id, u.uuid, u.first_name, u.last_name, u.email, u.role, u.password, u.group
`

const announcementFields = `
	a.id, a.uuid, a.name, a.content, a.date
`

const assignmentFields = `
	a.id, a.uuid, a.name, a.due_date, a.points, a.published
`

const discussionFields = `
	d.id, d.uuid, d.name, d.date,
	u.id, u.uuid, u.first_name, u.last_name, u.email, u.role, u.password, u.group
`

const quizFields = `
	q.id, q.uuid, q.name, q.due_date, q.points, q.published
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

	assignemnts, err := scanRowsToAssignments(rows)
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

	announcements, err := scanRowsToAnnouncements(rows)
	if err != nil {
		return nil, err
	}

	return announcements, nil
}

func (dao *CourseDao) getDiscussionsOfCourse(courseId int) ([]model.Discussion, error) {
	rows, err := dao.db.Query(
		fmt.Sprintf(`
			SELECT %s
			FROM discussions d
			JOIN users u on u.id = d.creator_id
			WHERE d.course_id = ?
			`,
			discussionFields,
		),
		courseId,
	)

	if err != nil {
		return nil, fmt.Errorf("Cannot query db: %w", err)
	}
	defer rows.Close()

	discussions, err := scanRowsToDiscussions(rows)
	if err != nil {
		return nil, err
	}

	return discussions, nil
}

func (dao *CourseDao) getQuizzesOfCourse(courseId int) ([]model.Quiz, error) {
	rows, err := dao.db.Query(
		fmt.Sprintf(`
			SELECT %s
			FROM quizzes q
			WHERE q.course_id = ?
			`,
			quizFields,
		),
		courseId,
	)

	if err != nil {
		return nil, fmt.Errorf("Cannot query db: %w", err)
	}
	defer rows.Close()

	quizzes, err := scanRowsToQuizzes(rows)
	if err != nil {
		return nil, err
	}

	return quizzes, nil
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

		discussions, err := dao.getDiscussionsOfCourse(courseId)
		if err != nil {
			return nil, fmt.Errorf("Cannot get discussions of course '%d': %w", courseId, err)
		}

		quizzes, err := dao.getQuizzesOfCourse(courseId)
		if err != nil {
			return nil, fmt.Errorf("Cannot get quizzes of course '%d': %w", courseId, err)
		}

		courses[index].Assignments = assignments
		courses[index].Announcements = announcements
		courses[index].Discussions = discussions
		courses[index].Quizzes = quizzes
	}

	return courses, nil
}

func scanRowsToAssignments(rows *sql.Rows) ([]model.Assignment, error) {
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

func scanRowsToAnnouncements(rows *sql.Rows) ([]model.Announcement, error) {
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

func scanRowsToDiscussions(rows *sql.Rows) ([]model.Discussion, error) {
	var discussions []model.Discussion

	for rows.Next() {
		var discussion model.Discussion
		var creator model.User
		var creatorGroup sql.NullString

		err := rows.Scan(
			&discussion.Id,
			&discussion.Uuid,
			&discussion.Name,
			&discussion.Date,
			&creator.Id,
			&creator.Uuid,
			&creator.FirstName,
			&creator.LastName,
			&creator.Email,
			&creator.Role,
			&creator.Password,
			&creatorGroup,
		)

		if err != nil {
			return nil, fmt.Errorf("Cannot scan discussion rows into model: %w", err)
		}

		utils.ConvertNullString(creatorGroup, &creator.Group)
		discussions = append(discussions, discussion)
	}

	return discussions, nil
}

func scanRowsToQuizzes(rows *sql.Rows) ([]model.Quiz, error) {
	var quizzes []model.Quiz

	for rows.Next() {
		var quiz model.Quiz

		err := rows.Scan(
			&quiz.Id,
			&quiz.Uuid,
			&quiz.Name,
			&quiz.DueDate,
			&quiz.Points,
			&quiz.Published,
		)

		if err != nil {
			return nil, fmt.Errorf("Cannot scan quiz rows into model: %w", err)
		}

		quizzes = append(quizzes, quiz)
	}

	return quizzes, nil
}

func scanRowsToCourses(rows *sql.Rows) ([]model.Course, error) {
	var courses []model.Course

	for rows.Next() {
		var course model.Course
		var teacher model.User
		var teacherGroup sql.NullString

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
			&teacherGroup,
		)

		if err != nil {
			return nil, fmt.Errorf("Cannot scan course row into model: %w", err)
		}

		utils.ConvertNullString(teacherGroup, &teacher.Group)
		courses = append(courses, course)
	}

	return courses, nil
}
