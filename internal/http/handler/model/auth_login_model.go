package model

type AuthLoginRequestModel struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthLoginUserResponseModel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AuthLoginResponseModel struct {
	Token string                     `json:"token"`
	User  AuthLoginUserResponseModel `json:"user"`
}
