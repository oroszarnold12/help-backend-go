package dao

import (
	"database/sql"
	"fmt"
	"help/model"
)

const assignmentGradeFields = `
	ag.id, ag.uuid, ag.grade,
	u.id, u.uuid, u.first_name, u.last_name, u.email, u.role, u.password,
	a.id, a.uuid, a.name, a.due_date, a.points, a.published
`

type AssignmentGradeDao struct {
	db *sql.DB
}

func NewAssignmentGradeDao(db *sql.DB) *AssignmentGradeDao {
	return &AssignmentGradeDao{db: db}
}

func (dao *AssignmentGradeDao) GetAssignmentGrades(courseId, submitterId int) ([]model.AssignmentGrade, error) {
	rows, err := dao.db.Query(
		fmt.Sprintf(`
			SELECT %s
			FROM assignment_grades ag
			JOIN assignments a ON a.id = ag.assignment_id 
			JOIN users u ON u.id = ag.submitter_id
			WHERE a.course_id = ? AND ag.submitter_id = ?
			`,
			assignmentGradeFields,
		),
		courseId,
		submitterId,
	)

	if err != nil {
		return nil, fmt.Errorf("Cannot query db: %w", err)
	}
	defer rows.Close()

	assignmentGrades, err := scanRowsToAssignmentGrades(rows)
	if err != nil {
		return nil, err
	}

	return assignmentGrades, nil
}

func scanRowsToAssignmentGrades(rows *sql.Rows) ([]model.AssignmentGrade, error) {
	var assignmentGrades []model.AssignmentGrade

	for rows.Next() {
		var assignmentGrade model.AssignmentGrade
		var submitter model.User
		var assignment model.Assignment

		err := rows.Scan(
			&assignmentGrade.Id,
			&assignmentGrade.Uuid,
			&assignmentGrade.Grade,
			&submitter.Id,
			&submitter.Uuid,
			&submitter.FirstName,
			&submitter.LastName,
			&submitter.Email,
			&submitter.Role,
			&submitter.Password,
			&assignment.Id,
			&assignment.Uuid,
			&assignment.Name,
			&assignment.DueDate,
			&assignment.Points,
			&assignment.Published,
		)

		if err != nil {
			return nil, fmt.Errorf("Cannot scan assignment grade rows into model: %w", err)
		}

		assignmentGrade.Submitter = submitter
		assignmentGrade.Assignment = assignment
		assignmentGrades = append(assignmentGrades, assignmentGrade)
	}

	return assignmentGrades, nil
}
