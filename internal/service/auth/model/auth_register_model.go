package model

import "time"

type RegisterUserRequestModel struct {
	Name     string
	Email    string
	Password string
}

type RegisterUserResponseModel struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
}
