package model

import "time"

type RegisterUserRequest struct {
	Name     string
	Email    string
	Password string
}

type RegisterUserResponse struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
}
