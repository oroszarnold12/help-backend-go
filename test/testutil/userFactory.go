package testutil

import (
	"help/model"

	"github.com/google/uuid"
)

func NewTestUser(overrides ...func(*model.User)) model.User {
	user := model.User{
		Id:        0,
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
