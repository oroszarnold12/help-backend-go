package model

import "help/dto"

type UserModel struct {
	Id        int
	Uuid      string
	FirstName string
	LastName  string
	Email     string
	Password  string
}

func (user UserModel) ToDto() dto.UserDto {
	return dto.UserDto{
		Id:        user.Id,
		Uuid:      user.Uuid,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
}
