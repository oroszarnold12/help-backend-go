package model

import (
	"github.com/google/uuid"
	"help/dto"
)

type UserModel struct {
	Id        int
	Uuid      uuid.UUID
	FirstName string
	LastName  string
	Email     string
	Password  string
}

func (user UserModel) ToDto() dto.UserDto {
	return dto.UserDto{
		Id:        user.Id,
		Uuid:      user.Uuid.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
}

func UserModelFromPostDto(dto dto.UserPostDto, password string) UserModel {
	return UserModel{
		Uuid:      uuid.New(),
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Email:     dto.Email,
		Password:  password,
	}
}
