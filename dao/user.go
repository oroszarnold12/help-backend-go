package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"help/errorsx"
	"help/model"
)

type UserDao struct {
	db *sql.DB
}

func NewUserDao(db *sql.DB) *UserDao {
	return &UserDao{db: db}
}

func (dao *UserDao) CreateUser(user model.UserModel) error {
	_, err := dao.db.Exec("INSERT INTO users (uuid, first_name, last_name, email, password) VALUES (?,?,?,?)", user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (dao *UserDao) GetUserByEmail(email string) (*model.UserModel, error) {
	row := dao.db.QueryRow("SELECT id, uuid, first_name, last_name, email FROM users where email = ?", email)

	userModel, err := scanRowIntoUserModel(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errorsx.NewNotFoundError("User", email)
		}

		return nil, fmt.Errorf("Cannot query db: %w", err)
	}

	return userModel, nil
}

func scanRowIntoUserModel(rows *sql.Row) (*model.UserModel, error) {
	userModel := &model.UserModel{}

	err := rows.Scan(
		&userModel.Id,
		&userModel.Uuid,
		&userModel.FirstName,
		&userModel.LastName,
		&userModel.Email,
	)

	if err != nil {
		return nil, fmt.Errorf("Cannot scan user row into model: %w", err)
	}

	return userModel, nil
}

func scanRowsIntoUserModel(rows *sql.Rows) (*model.UserModel, error) {
	userModel := &model.UserModel{}

	err := rows.Scan(
		&userModel.Id,
		&userModel.Uuid,
		&userModel.FirstName,
		&userModel.LastName,
		&userModel.Email,
	)

	if err != nil {
		return nil, fmt.Errorf("Cannot scan user row into model: %w", err)
	}

	return userModel, nil
}
