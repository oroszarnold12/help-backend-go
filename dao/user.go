package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"help/errorsx"
	"help/model"
	"help/utils"

	"github.com/google/uuid"
)

const userSelectFields = "id, uuid, first_name, last_name, email, role, password, `group`"

type UserLister interface {
	GetUsers() ([]model.User, error)
}

type UserDao struct {
	db *sql.DB
}

func NewUserDao(db *sql.DB) *UserDao {
	return &UserDao{db: db}
}

func (dao *UserDao) CreateUser(user model.User) error {
	_, err := dao.db.Exec("INSERT INTO users (uuid, first_name, last_name, email, role, password, `group`) VALUES (?,?,?,?,?,?,?)", user.Uuid, user.FirstName, user.LastName, user.Email, user.Role, user.Password, user.Group)
	if err != nil {
		return fmt.Errorf("Cannot exec statement: %w", err)
	}

	return nil
}

func (dao *UserDao) GetUserByEmail(email string) (*model.User, error) {
	row := dao.db.QueryRow(fmt.Sprintf("SELECT %s FROM users where email = ?", userSelectFields), email)

	user, err := scanRowToUser(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorsx.NewNotFoundError("User", email)
		}

		return nil, fmt.Errorf("Cannot query db: %w", err)
	}

	return user, nil
}

func (dao *UserDao) GetUserByUuid(uuid uuid.UUID) (*model.User, error) {
	row := dao.db.QueryRow(fmt.Sprintf("SELECT %s FROM users where uuid = ?", userSelectFields), uuid)

	user, err := scanRowToUser(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorsx.NewNotFoundError("User", uuid.String())
		}

		return nil, fmt.Errorf("Cannot query db: %w", err)
	}

	return user, nil
}

func (dao *UserDao) GetUsers() ([]model.User, error) {
	rows, err := dao.db.Query(fmt.Sprintf("SELECT %s FROM users", userSelectFields))

	if err != nil {
		return nil, fmt.Errorf("Cannot query db: %w", err)
	}
	defer rows.Close()

	users, err := scanRowsToUsers(rows)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func scanRowToUser(row *sql.Row) (*model.User, error) {
	var user model.User
	var userGroup sql.NullString

	err := row.Scan(
		&user.Id,
		&user.Uuid,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Role,
		&user.Password,
		&userGroup,
	)

	if err != nil {
		return nil, fmt.Errorf("Cannot scan user row into model: %w", err)
	}

	utils.ConvertNullString(userGroup, &user.Group)

	return &user, nil
}

func scanRowsToUsers(rows *sql.Rows) ([]model.User, error) {
	var users []model.User

	for rows.Next() {
		var user model.User
		var userGroup sql.NullString

		err := rows.Scan(
			&user.Id,
			&user.Uuid,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Role,
			&user.Password,
			&userGroup,
		)

		if err != nil {
			return nil, fmt.Errorf("Cannot scan user row into model: %w", err)
		}

		utils.ConvertNullString(userGroup, &user.Group)
		users = append(users, user)
	}

	return users, nil
}
