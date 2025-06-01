package dto

type UserDto struct {
	Id        int    `json:"id"`
	Uuid      string `json:"uuid"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type UserPostDto struct {
	FirstName string `json:"firstName" validate:"required,max=255"`
	LastName  string `json:"lastName" validate:"required,max=255"`
	Email     string `json:"email" validate:"required,email,max=255"`
	Password  string `json:"password" validate:"required,password,min=8,max=255"`
}
