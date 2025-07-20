package dao

import (
	"database/sql"
	"fmt"
	"help/model"
	"help/utils"
)

const quizGradeFields = `
	qg.id, qg.uuid, qg.grade,
	u.id, u.uuid, u.first_name, u.last_name, u.email, u.role, u.password, u.group,
	q.id, q.uuid, q.name, q.due_date, q.points, q.published
`

type QuizGradeDao struct {
	db *sql.DB
}

func NewQuizGradeDao(db *sql.DB) *QuizGradeDao {
	return &QuizGradeDao{db: db}
}

func (dao *QuizGradeDao) GetQuizGrades(courseId, submitterId int) ([]model.QuizGrade, error) {
	rows, err := dao.db.Query(
		fmt.Sprintf(`
			SELECT %s
			FROM quiz_grades qg
			JOIN quizzes q ON q.id = qg.quiz_id
			JOIN users u ON u.id = qg.submitter_id
			WHERE q.course_id = ? AND qg.submitter_id = ?
			`,
			quizGradeFields,
		),
		courseId,
		submitterId,
	)

	if err != nil {
		return nil, fmt.Errorf("Cannot query db: %w", err)
	}
	defer rows.Close()

	quizGrades, err := scanRowsToQuizGrades(rows)
	if err != nil {
		return nil, err
	}

	return quizGrades, nil
}

func scanRowsToQuizGrades(rows *sql.Rows) ([]model.QuizGrade, error) {
	var quizGrades []model.QuizGrade

	for rows.Next() {
		var quizGrade model.QuizGrade
		var submitter model.User
		var submitterGroup sql.NullString
		var quiz model.Quiz

		err := rows.Scan(
			&quizGrade.Id,
			&quizGrade.Uuid,
			&quizGrade.Grade,
			&submitter.Id,
			&submitter.Uuid,
			&submitter.FirstName,
			&submitter.LastName,
			&submitter.Email,
			&submitter.Role,
			&submitter.Password,
			&submitterGroup,
			&quiz.Id,
			&quiz.Uuid,
			&quiz.Name,
			&quiz.DueDate,
			&quiz.Points,
			&quiz.Published,
		)

		if err != nil {
			return nil, fmt.Errorf("Cannot scan quiz grade rows into model: %w", err)
		}

		utils.ConvertNullString(submitterGroup, &submitter.Group)
		quizGrade.Submitter = submitter
		quizGrade.Quiz = quiz
		quizGrades = append(quizGrades, quizGrade)
	}

	return quizGrades, nil
}
