package model

type LoginUserRequestModel struct {
	Email    string
	Password string
}

type LoginUserResponseModel struct {
	Token string
	ID    string
	Name  string
}
