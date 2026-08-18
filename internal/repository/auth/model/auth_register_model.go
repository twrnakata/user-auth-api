package model

import "time"

type CreateRegisterUserRequestModel struct {
	Name         string
	Email        string
	PasswordHash string
}

type CreateRegisterUserModel struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
}

type CreateRegisterUserDocumentModel struct {
	Name         string    `bson:"name"`
	Email        string    `bson:"email"`
	PasswordHash string    `bson:"passwordHash"`
	CreatedAt    time.Time `bson:"createdAt"`
}
