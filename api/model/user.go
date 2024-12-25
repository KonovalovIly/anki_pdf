package api_model

import database_models "github.com/KonovalovIly/anki_pdf/database/model"

type UserRegisterPayload struct {
	Login    string `json:"login" validate:"required,min=4,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=30"`
}

func (user *UserRegisterPayload) MapToDatabaseUser() *database_models.UserDto {
	return &database_models.UserDto{
		Login:    user.Login,
		Email:    user.Email,
		Password: user.Password,
	}
}
