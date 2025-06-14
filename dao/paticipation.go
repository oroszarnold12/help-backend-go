package dao

import (
	"database/sql"
	"fmt"
	"help/model"
)

const participationSelectFields = `
	p.id, p.uuid, p.show_on_dashboard, 
	u.id, u.uuid, u.first_name, u.last_name, u.email, u.role, u.password, 
	c.id, c.uuid, c.name, c.long_name, c.description
`

type ParticipationDao struct {
	db *sql.DB
}

func NewPariticipationDao(db *sql.DB) *ParticipationDao {
	return &ParticipationDao{db: db}
}

func (dao *ParticipationDao) GetParticipationsOfUser(userId int) ([]model.Participation, error) {
	rows, err := dao.db.Query(
		fmt.Sprintf(`
			SELECT %s 
			FROM participations p 
			JOIN users u ON u.id = p.user_id 
			JOIN courses c ON c.id = p.course_id 
			WHERE p.user_id = ?`,
			participationSelectFields,
		),
		userId,
	)

	if err != nil {
		return nil, fmt.Errorf("Cannot query db: %w", err)
	}
	defer rows.Close()

	participations, err := scanRowsToParticipations(rows)
	if err != nil {
		return nil, err
	}

	return participations, nil
}

func scanRowsToParticipations(rows *sql.Rows) ([]model.Participation, error) {
	var participations []model.Participation

	for rows.Next() {
		var participation model.Participation
		var course model.Course
		var user model.User
		err := rows.Scan(
			&participation.Id,
			&participation.Uuid,
			&participation.ShowOnDashboard,
			&user.Id,
			&user.Uuid,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Role,
			&user.Password,
			&course.Id,
			&course.Uuid,
			&course.Name,
			&course.LongName,
			&course.Description,
		)

		if err != nil {
			return nil, fmt.Errorf("Cannot scan participation row into model: %w", err)
		}

		participation.Course = course
		participation.User = user

		participations = append(participations, participation)
	}

	return participations, nil
}
