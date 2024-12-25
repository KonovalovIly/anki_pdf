package database_models

type UserDto struct {
	ID       int64  `json:"id"`
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"-"`
}
