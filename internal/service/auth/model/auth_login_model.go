package model

type LoginUserRequest struct {
	Email    string
	Password string
}

type LoginUserResponse struct {
	Token string
	ID    string
	Name  string
}
