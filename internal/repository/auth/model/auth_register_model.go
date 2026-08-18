package model

import "time"

type CreateRegisterUserRequest struct {
	Name         string
	Email        string
	PasswordHash string
}

type CreateRegisterUserResponse struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
}
