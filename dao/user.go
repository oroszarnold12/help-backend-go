package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"help/errorsx"
	"help/model"
)

const userSelectFields = "id, uuid, first_name, last_name, email, password"

type UserDao struct {
	db *sql.DB
}

func NewUserDao(db *sql.DB) *UserDao {
	return &UserDao{db: db}
}

func (dao *UserDao) CreateUser(user model.User) error {
	_, err := dao.db.Exec("INSERT INTO users (uuid, first_name, last_name, email, password) VALUES (?,?,?,?, ?)", user.Uuid, user.FirstName, user.LastName, user.Email, user.Password)
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

func (dao *UserDao) GetUserByUuid(uuid string) (*model.User, error) {
	row := dao.db.QueryRow(fmt.Sprintf("SELECT %s FROM users where uuid = ?", userSelectFields), uuid)

	user, err := scanRowToUser(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorsx.NewNotFoundError("User", uuid)
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

	err := row.Scan(
		&user.Id,
		&user.Uuid,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
	)

	if err != nil {
		return nil, fmt.Errorf("Cannot scan user row into model: %w", err)
	}

	return &user, nil
}

func scanRowsToUsers(rows *sql.Rows) ([]model.User, error) {
	var users []model.User

	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.Id,
			&user.Uuid,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Password,
		)

		if err != nil {
			return nil, fmt.Errorf("Cannot scan user row into model: %w", err)
		}

		users = append(users, user)
	}

	return users, nil
}
