package model

type AuthLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthLoginUserResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AuthLoginResponse struct {
	Token string                `json:"token"`
	User  AuthLoginUserResponse `json:"user"`
}
