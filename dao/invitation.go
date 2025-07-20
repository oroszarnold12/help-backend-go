package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"help/errorsx"
	"help/model"
	"help/utils"
	"strconv"
)

type InvitationDao struct {
	db *sql.DB
}

func NewInvitationDao(db *sql.DB) *InvitationDao {
	return &InvitationDao{db: db}
}

const invitationSelectionFields = `
	i.id, i.uuid, 
	u.id, u.uuid, u.first_name, u.last_name, u.email, u.role, u.password, u.group,
	c.id, c.uuid, c.name, c.long_name, c.description
`

func (dao *InvitationDao) GetInvitationsOfUser(userId int) ([]model.Invitation, error) {
	rows, err := dao.db.Query(
		fmt.Sprintf(`
			SELECT %s 
			FROM invitations i
			JOIN users u ON u.id = i.user_id
			JOIN courses c ON c.id = i.course_id
			WHERE u.id = ?`,
			invitationSelectionFields,
		),
		userId,
	)

	if err != nil {
		return nil, fmt.Errorf("Cannot query db: %w", err)
	}
	defer rows.Close()

	invitations, err := scanRowsToInvitations(rows)
	if err != nil {
		return nil, err
	}

	return invitations, nil
}

func (dao *InvitationDao) GetInvitationOfUser(userId int, invitationId int) (*model.Invitation, error) {
	row := dao.db.QueryRow(
		fmt.Sprintf(`
			SELECT %s 
			FROM invitations i
			JOIN users u ON u.id = i.user_id
			JOIN courses c ON c.id = i.course_id
			WHERE u.id = ? AND i.id = ?`,
			invitationSelectionFields,
		),
		userId,
		invitationId,
	)

	invitation, err := scanRowToInvitation(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorsx.NewNotFoundError("Invitation", strconv.Itoa(invitationId))
		}

		return nil, fmt.Errorf("Cannot query db: %w", err)
	}

	return invitation, nil
}

func (dao *InvitationDao) DeleteInvitation(invitationId int) error {
	_, err := dao.db.Exec("DELETE FROM invitations i WHERE i.id = ?", invitationId)
	if err != nil {
		return fmt.Errorf("Cannot exec statement: %w", err)
	}

	return nil
}

func scanRowToInvitation(row *sql.Row) (*model.Invitation, error) {
	var invitation model.Invitation
	var course model.Course
	var user model.User
	var userGroup sql.NullString

	err := row.Scan(
		&invitation.Id,
		&invitation.Uuid,
		&user.Id,
		&user.Uuid,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Role,
		&user.Password,
		&userGroup,
		&course.Id,
		&course.Uuid,
		&course.Name,
		&course.LongName,
		&course.Description,
	)

	if err != nil {
		return nil, fmt.Errorf("Cannot scan invitation row into model: %w", err)
	}

	utils.ConvertNullString(userGroup, &user.Group)
	invitation.User = user
	invitation.Course = course

	return &invitation, nil
}

func scanRowsToInvitations(rows *sql.Rows) ([]model.Invitation, error) {
	var invitations []model.Invitation

	for rows.Next() {
		var invitation model.Invitation
		var course model.Course
		var user model.User
		var userGroup sql.NullString

		err := rows.Scan(
			&invitation.Id,
			&invitation.Uuid,
			&user.Id,
			&user.Uuid,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Role,
			&user.Password,
			&userGroup,
			&course.Id,
			&course.Uuid,
			&course.Name,
			&course.LongName,
			&course.Description,
		)

		if err != nil {
			return nil, fmt.Errorf("Cannot scan invitation row into model: %w", err)
		}

		utils.ConvertNullString(userGroup, &user.Group)
		invitation.User = user
		invitation.Course = course

		invitations = append(invitations, invitation)
	}

	return invitations, nil
}
