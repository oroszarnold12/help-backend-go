package factory

import (
	"help/dto"
	"help/model"

	"github.com/google/uuid"
)

func NewTestUser(overrides ...func(*model.User)) model.User {
	user := model.User{
		Id:        1,
		Uuid:      uuid.New(),
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  "hashed-password",
		Role:      model.RoleStudent,
		Group:     "group",
	}

	for _, override := range overrides {
		override(&user)
	}
	return user
}

func NewTestUserPostDto(overrides ...func(*dto.UserPostDto)) dto.UserPostDto {
	userDto := dto.UserPostDto{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  "AdminAdmin1;",
		Role:      "ROLE_STUDENT",
		Group:     "group",
	}

	for _, override := range overrides {
		override(&userDto)
	}
	return userDto
}
