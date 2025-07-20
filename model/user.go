package model

import (
	"github.com/google/uuid"
	"help/dto"
)

type User struct {
	Id        int
	Uuid      uuid.UUID
	FirstName string
	LastName  string
	Email     string
	Password  string
	Role      Role
	Group     string
}

func (user User) ToDto() dto.UserGetDto {
	return dto.UserGetDto{
		Id:        user.Id,
		Uuid:      user.Uuid.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      string(user.Role),
		Group:     user.Group,
	}
}

func UserFromPostDto(dto dto.UserPostDto, password string) User {
	return User{
		Uuid:      uuid.New(),
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Email:     dto.Email,
		Password:  password,
		Role:      Role(dto.Role),
		Group:     dto.Group,
	}
}
