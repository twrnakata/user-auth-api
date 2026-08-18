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

type CreateRegisterUserDocument struct {
	Name         string    `bson:"name"`
	Email        string    `bson:"email"`
	PasswordHash string    `bson:"passwordHash"`
	CreatedAt    time.Time `bson:"createdAt"`
}
